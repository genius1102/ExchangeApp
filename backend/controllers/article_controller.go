package controllers

import (
	"errors"
	"net/http"

	"exchangeapp/backend/global"
	"exchangeapp/backend/models"
	"exchangeapp/backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateArticle 创建文章
// Controller 职责：绑定参数 → 调用 Service → 返回 HTTP 响应
func CreateArticle(ctx *gin.Context) {
	var article models.Article

	if err := ctx.ShouldBindBodyWithJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.ArticleSvc.CreateArticle(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

// GetArticles 分页获取文章列表
func GetArticles(ctx *gin.Context) {
	page, size := utils.ParsePagination(ctx.Query("page"), ctx.Query("size"))

	resp, err := global.ArticleSvc.GetArticles(page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetArticleByID 根据ID获取文章详情
func GetArticleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	article, err := global.ArticleSvc.GetArticleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, article)
}
