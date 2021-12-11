package app

import (
	"github.com/dleonsal/beers-api/src/configs"
	"github.com/dleonsal/beers-api/src/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func StartApplication() {
	logger.Log.Info("Starting Application")
	router := gin.Default()
	config := configs.NewConfig()
	handlers := wireDependencies(config)

	mapRoutes(router, handlers)

	router.Run(":" + config.Port)
}
