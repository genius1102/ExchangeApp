package controllers

import (
	"exchangeapp/backend/global"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// toggleLikeScript 使用Lua脚本保证「判断+去重+计数」的原子性
// 返回 1 表示点赞成功，0 表示取消点赞成功
const toggleLikeScript = `
local userKey = KEYS[1]
local likeKey = KEYS[2]
local username = ARGV[1]

local isMember = redis.call('SISMEMBER', userKey, username)
if isMember == 1 then
    redis.call('SREM', userKey, username)
    redis.call('DECR', likeKey)
    return 0
else
    redis.call('SADD', userKey, username)
    redis.call('INCR', likeKey)
    return 1
end
`

func ToggleLike(ctx *gin.Context) {
	articleID := ctx.Param("id")
	username, _ := ctx.Get("username")

	userKey := "article:" + articleID + ":liked_users"
	likeKey := "article:" + articleID + ":likes"

	result, err := global.RedisDB.Eval(toggleLikeScript, []string{userKey, likeKey}, username.(string)).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	liked := result.(int64) == 1
	ctx.JSON(http.StatusOK, gin.H{"liked": liked, "message": "success"})
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
