package controllers

import (
	"net/http"

	"exchangeapp/backend/global"
	"exchangeapp/backend/models"

	"github.com/gin-gonic/gin"
)

// CreateExchangeRate 创建汇率记录
func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.ExchangeRateSvc.CreateExchangeRate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, exchangeRate)
}

// GetExchangeRates 获取全部汇率记录
func GetExchangeRates(ctx *gin.Context) {
	exchangeRates, err := global.ExchangeRateSvc.GetExchangeRates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, exchangeRates)
}
