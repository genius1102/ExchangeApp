package controllers

import (
	"exchangeapp/backend/global"
	"exchangeapp/backend/models"
	"exchangeapp/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user models.User

	// 绑定 JSON 数据到 user 结构体
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"err":err.Error(),
		})
		return
	}

	// 密码加密
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}

	user.Password = hash

	// JWT生成
	token,err := utils.GenerateJWT(user.UserName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}


	// 数据库迁移和用户创建
	if err := global.DB.AutoMigrate(&user); err != nil	{
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}


	// 创建用户记录
	if err := global.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}


	// 返回 JWT
	ctx.JSON(http.StatusOK,gin.H{
		"token":token,
	})

}

func Login(ctx *gin.Context) {
	var input struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"err":err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User 
	if err := global.DB.Where("user_name = ?",input.UserName).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized,gin.H{
			"err":"user not found",
		})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized,gin.H{
			"err":"password incorrect",
		})
		return
	}

	// 生成 JWT
	token, err := utils.GenerateJWT(user.UserName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}

	// 返回 JWT
	ctx.JSON(http.StatusOK,gin.H{
		"token":token,
	})
}
