package services

import (
	"fmt"
	"strconv"

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

type LikeService struct {
	redis *redis.Client
}

func NewLikeService(redis *redis.Client) *LikeService {
	return &LikeService{redis: redis}
}

// ToggleLike 切换点赞状态，返回 (是否已赞, error)
func (s *LikeService) ToggleLike(articleID, username string) (bool, error) {
	userKey := fmt.Sprintf("article:%s:liked_users", articleID)
	likeKey := fmt.Sprintf("article:%s:likes", articleID)

	result, err := s.redis.Eval(toggleLikeScript, []string{userKey, likeKey}, username).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

// GetArticleLikes 获取点赞数和当前用户的点赞状态
func (s *LikeService) GetArticleLikes(articleID, username string) (likes int, liked bool, err error) {
	likeKey := fmt.Sprintf("article:%s:likes", articleID)
	userKey := fmt.Sprintf("article:%s:liked_users", articleID)

	likesStr, err := s.redis.Get(likeKey).Result()
	if err == redis.Nil {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}

	likes, _ = strconv.Atoi(likesStr)

	isMember, err := s.redis.SIsMember(userKey, username).Result()
	if err != nil {
		return 0, false, err
	}

	return likes, isMember, nil
}
