package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/backend/global"
	"exchangeapp/backend/models"
	"exchangeapp/backend/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

const (
	defaultPageSize = 10
	maxChangePages  = 3
)

func CreateArticle(ctx *gin.Context) {

	var article models.Article

	if err := ctx.ShouldBindBodyWithJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除缓存,防止用户在缓存有效期内浏览文章列表看不到新增的文章
	for i := 0; i < maxChangePages; i++ {
		cacheKey := fmt.Sprintf("articles:page:%d:size:%d", i, defaultPageSize)
		if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	ctx.JSON(http.StatusCreated, article)

}

// 获取文章列表
func GetArticles(ctx *gin.Context) {

	page, size := utils.ParsePagination(ctx.Query("page"), ctx.Query("size"))
	if page <= maxChangePages {
		getArticlesWithCache(ctx, page, size)
	} else {
		getArticlesWithDB(ctx, page, size)
	}

}

func getArticlesWithCache(ctx *gin.Context, page, size int) {
	cacheKey := fmt.Sprintf("articles:page:%d:size:%d", page, size)
	cachedData, err := global.RedisDB.Get(cacheKey).Result()

	// 缓存命中
	if err == nil {
		var resp models.PageResponse
		if err := json.Unmarshal([]byte(cachedData), &resp); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	// 其他错误
	if err != redis.Nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 缓存未命中，查 MySQL
	resp, err := queryArticlesFromDB(page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)

	// 缓存文章列表
	cacheData, err := json.Marshal(resp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := global.RedisDB.Set(cacheKey, cacheData, 10*time.Minute).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

// getArticlesWithDB 第三页开始从数据库查询文章列表
func getArticlesWithDB(ctx *gin.Context, page, size int) {
	resp, err := queryArticlesFromDB(page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, resp)
}

// queryArticlesFromDB 查询文章列表
func queryArticlesFromDB(page, size int) (*models.PageResponse, error) {
	var articles []models.Article
	var total int64

	global.DB.Model(&models.Article{}).Count(&total)

	if err := global.DB.Limit(size).Offset((page - 1) * size).Find(&articles).Error; err != nil {
		return nil, err
	}

	return &models.PageResponse{
		Data:  articles,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var article models.Article

	if err := global.DB.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}
