package app

import (
	"github.com/gin-gonic/gin"
)

func mapRoutes(router *gin.Engine, handlers *handlerContainer) {
	router.GET("/beers", handlers.beerHandler.HandleList)
	router.GET("/beers/:beer_id", handlers.beerHandler.HandleGetByID)
	router.GET("/beers/:beer_id/boxprice", handlers.beerHandler.HandleGetBoxPrice)
	router.POST("/beers", handlers.beerHandler.HandleCreate)

}
