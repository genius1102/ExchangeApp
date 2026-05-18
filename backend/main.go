package main

import (
	"exchangeapp/backend/config"
	"exchangeapp/backend/router"
	"fmt"
)

func main() {
	config.InitConfig()
	fmt.Println(config.AppConfig.App.Name)

	router := router.SetUpRouter()

	router.Run(":" + fmt.Sprintf("%d", config.AppConfig.App.Port))
}