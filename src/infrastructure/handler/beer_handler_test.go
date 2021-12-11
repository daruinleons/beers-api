package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dleonsal/beers-api/src/core/contracts"
	"github.com/dleonsal/beers-api/src/core/domain/entities"
	"github.com/dleonsal/beers-api/src/errors"
	"github.com/dleonsal/beers-api/src/infrastructure/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_HandleList_WhenBeerServiceFail_ThenReturnErrorAndStatusCode(t *testing.T) {
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers", nil, nil, "")
	expectedError := errors.NewInternalServerError("some error")
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("ListBeers").Return(nil, expectedError)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleList(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleList_WhenProcessIsExecutedCorrectly_ThenReturnBeersAndStatusCode200(t *testing.T) {
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers", nil, nil, "")
	expectedBeer := givenBeer()
	expectedBeers := []entities.Beer{*expectedBeer}
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("ListBeers").Return(expectedBeers, nil)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleList(ctx)

	beers := new([]entities.Beer)
	json.Unmarshal(recorder.Body.Bytes(), beers)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, expectedBeers, *beers)
}

func Test_HandleGetByID_WhenParamBeerIDIsInvalid_ThenReturnErrorAndStatusCode(t *testing.T) {
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers/:beer_id",
		[]gin.Param{{Key: "beer_id", Value: "invalid"}}, nil, "")
	expectedError := errors.NewBadRequestError("id should be a number")
	handler := handler.NewBeerHandler(nil)

	handler.HandleGetByID(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleGetByID_WhenBeerServiceFail_ThenReturnErrorAndStatusCode(t *testing.T) {
	id := int64(1)
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers/:beer_id",
		[]gin.Param{{Key: "beer_id", Value: fmt.Sprint(id)}}, nil, "")
	expectedError := errors.NewInternalServerError("some error")
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("GetBeerByID", id).Return(nil, expectedError)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleGetByID(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleGetByID_WhenProcessIsExecutedCorrectly_ThenReturnBeerAndStatusCode200(t *testing.T) {
	expectedBeer := givenBeer()
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers/:beer_id",
		[]gin.Param{{Key: "beer_id", Value: fmt.Sprint(expectedBeer.Id)}}, nil, "")
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("GetBeerByID", expectedBeer.Id).Return(expectedBeer, nil)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleGetByID(ctx)

	beer := new(entities.Beer)
	json.Unmarshal(recorder.Body.Bytes(), beer)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, expectedBeer, beer)
}

func Test_HandleGetBoxPrice_WhenParamBeerIDIsInvalid_ThenReturnErrorAndStatusCode(t *testing.T) {
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers/:beer_id/boxprice",
		[]gin.Param{{Key: "beer_id", Value: "invalid"}}, nil, "")
	expectedError := errors.NewBadRequestError("id should be a number")
	handler := handler.NewBeerHandler(nil)

	handler.HandleGetBoxPrice(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleGetBoxPrice_WhenParamQuantityIsInvalid_ThenReturnErrorAndStatusCode(t *testing.T) {
	queryParams := url.Values{"quantity": {"invalid"}}
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers/:beer_id/boxprice",
		[]gin.Param{{Key: "beer_id", Value: "1"}}, &queryParams, "")
	expectedError := errors.NewBadRequestError("quantity should be a positive number")
	handler := handler.NewBeerHandler(nil)

	handler.HandleGetBoxPrice(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleGetBoxPrice_WhenBeerServiceFail_ThenReturnErrorAndStatusCode(t *testing.T) {
	id := int64(1)
	newCurrency := "USD"
	quantity := uint64(10)
	queryParams := url.Values{"currency": {newCurrency}, "quantity": {fmt.Sprint(quantity)}}
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers/:beer_id/boxprice",
		[]gin.Param{{Key: "beer_id", Value: fmt.Sprint(id)}}, &queryParams, "")
	expectedError := errors.NewInternalServerError("some error")
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("GetBoxPrice", id, newCurrency, quantity).Return(0.0, expectedError)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleGetBoxPrice(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleGetBoxPrice_WhenProcessIsExecutedCorrectly_ThenReturnErrorAndStatusCode(t *testing.T) {
	id := int64(1)
	newCurrency := "USD"
	quantity := uint64(10)
	queryParams := url.Values{"currency": {newCurrency}, "quantity": {fmt.Sprint(quantity)}}
	expectedTotalPrice := 10.0
	expectedResponse := contracts.BoxPriceResponse{
		TotalPrice: expectedTotalPrice,
	}
	ctx, recorder := givenContextAndRecorder(http.MethodGet, "/beers/:beer_id/boxprice",
		[]gin.Param{{Key: "beer_id", Value: fmt.Sprint(id)}}, &queryParams, "")
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("GetBoxPrice", id, newCurrency, quantity).Return(expectedTotalPrice, nil)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleGetBoxPrice(ctx)

	boxPriceResponse := new(contracts.BoxPriceResponse)
	json.Unmarshal(recorder.Body.Bytes(), boxPriceResponse)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, expectedResponse, *boxPriceResponse)
}

func Test_HandleCreate_WhenBodyIsInvalid_ThenReturnErrorAndStatusCode(t *testing.T) {
	ctx, recorder := givenContextAndRecorder(http.MethodPost, "/beers",
		nil, nil, "{,}")
	expectedError := errors.NewBadRequestError("invalid json body")
	handler := handler.NewBeerHandler(nil)

	handler.HandleCreate(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleCreate_WhenBeerServiceFail_ThenReturnErrorAndStatusCode(t *testing.T) {
	beer := givenBeer()
	bodyBytes, _ := json.Marshal(beer)
	ctx, recorder := givenContextAndRecorder(http.MethodPost, "/beers",
		nil, nil, string(bodyBytes))
	expectedError := errors.NewInternalServerError("some error")
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("CreateBeer", *beer).Return(expectedError)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleCreate(ctx)

	assert.Equal(t, expectedError.Status, recorder.Code)
	assert.Equal(t, expectedError, getRestError(recorder.Body.Bytes()))
}

func Test_HandleCreate_WhenProcessIsExecutedCorrectly_ThenReturnErrorAndStatusCode(t *testing.T) {
	beer := givenBeer()
	bodyBytes, _ := json.Marshal(beer)
	ctx, recorder := givenContextAndRecorder(http.MethodPost, "/beers",
		nil, nil, string(bodyBytes))
	mockBeerService := new(handler.MockBeerService)
	mockBeerService.On("CreateBeer", *beer).Return(nil)
	handler := handler.NewBeerHandler(mockBeerService)

	handler.HandleCreate(ctx)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "\"Beer created\"", recorder.Body.String())
}

func givenContextAndRecorder(method, url string, params []gin.Param, queryParams *url.Values, body string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	if params != nil {
		ctx.Params = params
	}

	if queryParams != nil {
		url = url + "?" + queryParams.Encode()
	}

	ctx.Request, _ = http.NewRequest(method, url, strings.NewReader(body))

	return ctx, recorder
}

func getRestError(bodyBytes []byte) *errors.RestError {
	restError := new(errors.RestError)
	json.Unmarshal(bodyBytes, &restError)
	return restError
}

func givenBeer() *entities.Beer {
	return &entities.Beer{
		Id:       1,
		Name:     "Pilsen",
		Brewery:  "Bavaria",
		Country:  "Colombia",
		Price:    2500,
		Currency: "COP",
	}
}
