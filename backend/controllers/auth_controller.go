package controllers

import (
	"errors"
	"net/http"

	"exchangeapp/backend/global"
	"exchangeapp/backend/models"
	"exchangeapp/backend/utils"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := global.AuthSvc.Register(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Login 用户登录
func Login(ctx *gin.Context) {
	var input struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := global.AuthSvc.Login(input.UserName, input.Password)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) || errors.Is(err, utils.ErrInvalidPassword) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
