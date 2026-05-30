package config

import (
	"log"
	"os"
	"strconv"

	"exchangeapp/backend/utils"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name      string `yaml:"name"`
		Port      int    `yaml:"port"`
		JWTSecret string `yaml:"jwt_secret"`
	}
	Database struct {
		Dsn string `yaml:"dsn"`
		MaxIdleConns int `yaml:"MaxIdleConns"`
		MaxOpenConns int `yaml:"MaxOpenConns"`
	}
	Redis struct {
		Addr string `yaml:"addr"`
		Password string `yaml:"password"`
		DB int `yaml:"db"`
	}
}

var AppConfig *Config

func InitConfig(){
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v",err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("error unmarshalling config struct: %v",err)
	}

	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		AppConfig.Database.Dsn = dsn
	}
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		AppConfig.Redis.Addr = addr
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		AppConfig.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		dbInt, err := strconv.Atoi(db)
		if err == nil {
			AppConfig.Redis.DB = dbInt
		}
	}
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		AppConfig.App.JWTSecret = secret
	}

	utils.JWTSecret = AppConfig.App.JWTSecret

	InitDB()
	InitRedis()

}