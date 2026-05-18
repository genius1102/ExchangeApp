package global

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	RedisDB *redis.Client
)
