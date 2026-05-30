package controllers

import (
	"net/http"

	"exchangeapp/backend/global"

	"github.com/gin-gonic/gin"
)

// ToggleLike 切换点赞状态
func ToggleLike(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username, _ := ctx.Get("username")

	liked, err := global.LikeSvc.ToggleLike(articleID, username.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"liked": liked, "message": "success"})
}

// GetArticleLikes 获取点赞信息
func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username, _ := ctx.Get("username")

	likes, liked, err := global.LikeSvc.GetArticleLikes(articleID, username.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"likes": likes, "liked": liked})
}
