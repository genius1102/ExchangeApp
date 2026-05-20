package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/backend/global"
	"exchangeapp/backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var cacheKey = "articles"

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
	if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, article)

}

func GetArticles(ctx *gin.Context) {

	cachedData, err := global.RedisDB.Get(cacheKey).Result()

	// 1、缓存未命中
	if err == redis.Nil {
		var articles []models.Article

		if err := global.DB.Find(&articles).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "No articles found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// 返回文章列表
		ctx.JSON(http.StatusOK, articles)

		// 缓存文章列表
		articlesJSON, err := json.Marshal(articles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 缓存文章列表 用Set方法
		if err := global.RedisDB.Set(cacheKey, articlesJSON, 10*time.Minute).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
			return
		}

	} else if err != nil {
		// 2、其他错误
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		// 3、缓存命中
		var articles []models.Article
		if err := json.Unmarshal([]byte(cachedData),&articles); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, articles)
	}

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
