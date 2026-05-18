package middlewares

import (
	"exchangeapp/backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			ctx.Abort() // 中止后续中间件处理
			return
		}

		username, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}
		ctx.Set("username", username) // 存储用户名到上下文 Set方法可以在上下文中创建新的键值对在上下文中传递
		ctx.Next()                    //继续处理后续中间件

	}

}
