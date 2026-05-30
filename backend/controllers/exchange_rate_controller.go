package controllers

import (
	"net/http"
	"time"

	"exchangeapp/backend/global"
	"exchangeapp/backend/models"

	"github.com/gin-gonic/gin"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	// 绑定 JSON 数据到 exchangeRate 结构体
	if err := ctx.ShouldBindJSON(&exchangeRate) ; err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchangeRate.Date = time.Now()

	if err := global.DB.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, exchangeRate)
}

func GetExchangeRates(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate

	// 从数据库查询所有汇率记录
	if err := global.DB.Find(&exchangeRates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回查询结果
	ctx.JSON(http.StatusOK, exchangeRates)
}