package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
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

	InitDB()
	InitRedis()

}