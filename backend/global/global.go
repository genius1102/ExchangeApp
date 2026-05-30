package global

import (
	"exchangeapp/backend/services"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	RedisDB        *redis.Client
	ArticleSvc     *services.ArticleService
	AuthSvc        *services.AuthService
	ExchangeRateSvc *services.ExchangeRateService
	LikeSvc        *services.LikeService
)
