package config

import (
	"log"
	"time"

	"exchangeapp/backend/global"
	"exchangeapp/backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(100 * time.Second)

	err = db.AutoMigrate(&models.User{}, &models.Article{}, &models.ExchangeRate{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	global.DB = db
}
