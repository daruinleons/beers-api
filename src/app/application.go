package app

import (
	"github.com/dleonsal/beers-api/src/configs"
	"github.com/gin-gonic/gin"
)

func StartApplication() {
	router := gin.Default()
	config := configs.NewConfig()
	handlers := wireDependencies(config)

	mapRoutes(router, handlers)

	router.Run(":" + config.Port)
}
