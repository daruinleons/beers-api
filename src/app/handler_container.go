package app

import "github.com/gin-gonic/gin"

type beerHandler interface {
	HandleList(c *gin.Context)
	HandleGetByID(c *gin.Context)
	HandleGetBoxPrice(c *gin.Context)
	HandleCreate(c *gin.Context)
}

type handlerContainer struct {
	beerHandler beerHandler
}

func newHandlerContainer(beerHandler beerHandler) *handlerContainer {
	return &handlerContainer{
		beerHandler: beerHandler,
	}
}
