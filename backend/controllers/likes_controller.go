package controllers

import (
	"exchangeapp/backend/global"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func ToggleLike(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username, _ := ctx.Get("username")

	userKey := "article:" + articleID + ":liked_users"
	likeKey := "article:" + articleID + ":likes"

	isMember, err := global.RedisDB.SIsMember(userKey, username.(string)).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isMember {
		global.RedisDB.SRem(userKey, username.(string))
		global.RedisDB.Decr(likeKey)
		ctx.JSON(http.StatusOK, gin.H{"liked": false, "message": "successfully unliked the article!"})
		return
	}

	global.RedisDB.SAdd(userKey, username.(string))
	global.RedisDB.Incr(likeKey)
	ctx.JSON(http.StatusOK, gin.H{"liked": true, "message": "successfully liked the article!"})
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username, _ := ctx.Get("username")

	likeKey := "article:" + articleID + ":likes"
	userKey := "article:" + articleID + ":liked_users"

	likesStr, err := global.RedisDB.Get(likeKey).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusOK, gin.H{"likes": 0, "liked": false})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	likes, _ := strconv.Atoi(likesStr)

	isMember, err := global.RedisDB.SIsMember(userKey, username.(string)).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"likes": likes, "liked": isMember})
}
