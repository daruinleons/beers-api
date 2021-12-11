package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dleonsal/beers-api/src/core/contracts"
	"github.com/dleonsal/beers-api/src/core/domain/entities"
	"github.com/dleonsal/beers-api/src/errors"
	"github.com/dleonsal/beers-api/src/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

type BeerService interface {
	ListBeers() ([]entities.Beer, *errors.RestError)
	GetBeerByID(beerID int64) (*entities.Beer, *errors.RestError)
	GetBoxPrice(beerID int64, newCurrency string, quantity uint64) (float64, *errors.RestError)
	CreateBeer(beer entities.Beer) *errors.RestError
}

type beerHandler struct {
	beerService BeerService
}

func NewBeerHandler(beenService BeerService) *beerHandler {
	return &beerHandler{
		beerService: beenService,
	}
}

func (h *beerHandler) HandleList(c *gin.Context) {
	beers, restErr := h.beerService.ListBeers()
	if restErr != nil {
		c.JSON(restErr.Status, restErr)

		return
	}

	c.JSON(http.StatusOK, beers)
}

func (h *beerHandler) HandleGetByID(c *gin.Context) {
	beerID, err := strconv.ParseInt(c.Param("beer_id"), 10, 64)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to parse param beer id to int64: %s", err))
		restErr := errors.NewBadRequestError("id should be a number")
		c.JSON(restErr.Status, restErr)

		return
	}

	beer, restErr := h.beerService.GetBeerByID(beerID)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)

		return
	}

	c.JSON(http.StatusOK, beer)
}

func (h *beerHandler) HandleGetBoxPrice(c *gin.Context) {
	beerID, err := strconv.ParseInt(c.Param("beer_id"), 10, 64)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to parse param beer id to int64: %s", err))
		restErr := errors.NewBadRequestError("id should be a number")
		c.JSON(restErr.Status, restErr)

		return
	}

	currency := c.Query("currency")

	quantity, err := strconv.ParseUint(c.Query("quantity"), 10, 64)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to parse beer param quantity to uint64: %s", err))
		restErr := errors.NewBadRequestError("quantity should be a positive number")
		c.JSON(restErr.Status, restErr)

		return
	}

	totalPrice, restErr := h.beerService.GetBoxPrice(beerID, currency, quantity)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)

		return
	}

	c.JSON(http.StatusOK, contracts.BoxPriceResponse{
		TotalPrice: totalPrice,
	})

}

func (h *beerHandler) HandleCreate(c *gin.Context) {
	var request entities.Beer

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to bind request body: %s", err))
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)

		return
	}

	restErr := h.beerService.CreateBeer(request)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)

		return
	}

	c.JSON(http.StatusCreated, "Beer created")
}
