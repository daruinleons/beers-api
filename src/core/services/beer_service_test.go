package services_test

import (
	"fmt"
	"github.com/dleonsal/beers-api/src/core/domain/entities"
	"github.com/dleonsal/beers-api/src/core/services"
	"github.com/dleonsal/beers-api/src/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ListBeers_WhenRepositoryFail_ThenReturnError(t *testing.T) {
	expectedError := errors.NewInternalServerError("some error")
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("List").Return(nil, expectedError)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	beers, err := beerService.ListBeers()

	assert.Nil(t, beers)
	assert.Equal(t, expectedError, err)
	mockBeerRepository.AssertExpectations(t)

}

func Test_ListBeers_WhenProcessIsExecutedSuccessfully_ThenReturnBeersList(t *testing.T) {
	expectBeer := givenBeer()
	expectedBeers := []entities.Beer{*expectBeer}
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("List").Return(expectedBeers, nil)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	beers, err := beerService.ListBeers()

	assert.Equal(t, expectedBeers, beers)
	assert.Nil(t, err)
	mockBeerRepository.AssertExpectations(t)

}

func Test_GetBeerByID_WhenRepositoryFail_ThenReturnError(t *testing.T) {
	id := int64(1)
	expectedError := errors.NewInternalServerError("some error")
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("GetByID", id).Return(nil, expectedError)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	beer, err := beerService.GetBeerByID(id)

	assert.Nil(t, beer)
	assert.Equal(t, expectedError, err)
	mockBeerRepository.AssertExpectations(t)

}

func Test_GetBeerByID_WhenProcessIsExecutedSuccessfully_ThenReturnBeer(t *testing.T) {
	id := int64(1)
	expectBeer := givenBeer()
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("GetByID", id).Return(expectBeer, nil)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	beer, err := beerService.GetBeerByID(id)

	assert.Equal(t, expectBeer, beer)
	assert.Nil(t, err)
	mockBeerRepository.AssertExpectations(t)

}

func Test_GetBoxPrice_WhenNewCurrencyIsInvalid_ThenReturnError(t *testing.T) {
	id := int64(1)
	newCurrency := ""
	quantity := uint64(10)
	expectTotalPrice := 0.0
	expectedError := errors.NewBadRequestError("currency must not be empty")
	beerService := services.NewBeerService(nil, nil)

	totalPrice, err := beerService.GetBoxPrice(id, newCurrency, quantity)

	assert.Equal(t, expectTotalPrice, totalPrice)
	assert.Equal(t, expectedError, err)
}

func Test_GetBoxPrice_WhenGetBeerRepositoryFail_ThenReturnError(t *testing.T) {
	id := int64(1)
	newCurrency := "USD"
	quantity := uint64(10)
	expectedTotalPrice := 0.0
	expectedError := errors.NewInternalServerError("some error")
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("GetByID", id).Return(nil, expectedError)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	totalPrice, err := beerService.GetBoxPrice(id, newCurrency, quantity)

	assert.Equal(t, expectedTotalPrice, totalPrice)
	assert.Equal(t, expectedError, err)
	mockBeerRepository.AssertExpectations(t)
}

func Test_GetBoxPrice_WhenNewCurrencyIsEqualToCurrentCurrency_ThenReturnTotalPriceWithoutCallConverter(t *testing.T) {
	id := int64(1)
	newCurrency := "COP"
	quantity := uint64(10)
	expectBeer := givenBeer()
	expectedTotalPrice := 25000.0
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("GetByID", id).Return(expectBeer, nil)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	totalPrice, err := beerService.GetBoxPrice(id, newCurrency, quantity)

	assert.Equal(t, expectedTotalPrice, totalPrice)
	assert.Nil(t, err)
	mockBeerRepository.AssertExpectations(t)
}

func Test_GetBoxPrice_WhenCurrencyConverterClientFail_ThenReturnError(t *testing.T) {
	id := int64(1)
	newCurrency := "USD"
	quantity := uint64(10)
	expectBeer := givenBeer()
	expectedTotalPrice := 0.0
	expectedError := errors.NewInternalServerError("some error")
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("GetByID", id).Return(expectBeer, nil)
	mockCurrencyConverterClient := new(services.MockCurrencyConverterClient)
	mockCurrencyConverterClient.On("ConvertValueToNewCurrency", expectBeer.Currency, newCurrency, expectBeer.Price).
		Return(expectedTotalPrice, expectedError)
	beerService := services.NewBeerService(mockBeerRepository, mockCurrencyConverterClient)

	totalPrice, err := beerService.GetBoxPrice(id, newCurrency, quantity)

	assert.Equal(t, expectedTotalPrice, totalPrice)
	assert.Equal(t, expectedError, err)
	mockBeerRepository.AssertExpectations(t)
	mockCurrencyConverterClient.AssertExpectations(t)
}

func Test_GetBoxPrice_WhenProcessIsExecutedSuccessfully_ThenReturnTotalPrice(t *testing.T) {
	id := int64(1)
	newCurrency := "USD"
	quantity := uint64(10)
	expectBeer := givenBeer()
	expectedTotalPrice := 6.0
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("GetByID", id).Return(expectBeer, nil)
	mockCurrencyConverterClient := new(services.MockCurrencyConverterClient)
	mockCurrencyConverterClient.On("ConvertValueToNewCurrency", expectBeer.Currency, newCurrency, expectBeer.Price).
		Return(0.6, nil)
	beerService := services.NewBeerService(mockBeerRepository, mockCurrencyConverterClient)

	totalPrice, err := beerService.GetBoxPrice(id, newCurrency, quantity)

	assert.Equal(t, expectedTotalPrice, totalPrice)
	assert.Nil(t, err)
	mockBeerRepository.AssertExpectations(t)
	mockCurrencyConverterClient.AssertExpectations(t)
}

func Test_GetBoxPrice_WhenProcessIsExecutedSuccessfullyWithQuantityZero_ThenReturnTotalPrice(t *testing.T) {
	id := int64(1)
	newCurrency := "USD"
	quantity := uint64(0)
	expectBeer := givenBeer()
	expectedTotalPrice := 3.5999999999999996
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("GetByID", id).Return(expectBeer, nil)
	mockCurrencyConverterClient := new(services.MockCurrencyConverterClient)
	mockCurrencyConverterClient.On("ConvertValueToNewCurrency", expectBeer.Currency, newCurrency, expectBeer.Price).
		Return(0.6, nil)
	beerService := services.NewBeerService(mockBeerRepository, mockCurrencyConverterClient)

	totalPrice, err := beerService.GetBoxPrice(id, newCurrency, quantity)

	assert.Equal(t, expectedTotalPrice, totalPrice)
	assert.Nil(t, err)
	mockBeerRepository.AssertExpectations(t)
	mockCurrencyConverterClient.AssertExpectations(t)
}

func Test_CreateBeer_WhenBeerValidateFail_ThenReturnError(t *testing.T) {
	beer := entities.Beer{
		Id: 0,
	}
	expectedError := errors.NewBadRequestError(fmt.Sprintf("invalid Id: %d", beer.Id))
	beerService := services.NewBeerService(nil, nil)

	err := beerService.CreateBeer(beer)

	assert.Equal(t, expectedError, err)
}

func Test_CreateBeer_WhenRepositoryFail_ThenReturnError(t *testing.T) {
	beer := givenBeer()
	expectedError := errors.NewInternalServerError("some error")
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("Save", *beer).Return(expectedError)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	err := beerService.CreateBeer(*beer)

	assert.Equal(t, expectedError, err)
	mockBeerRepository.AssertExpectations(t)

}

func Test_CreateBeer_WhenProcessIsExecutedSuccessfully_ThenReturnNil(t *testing.T) {
	beer := givenBeer()
	mockBeerRepository := new(services.MockBeerRepository)
	mockBeerRepository.On("Save", *beer).Return(nil)
	beerService := services.NewBeerService(mockBeerRepository, nil)

	err := beerService.CreateBeer(*beer)

	assert.Nil(t, err)
	mockBeerRepository.AssertExpectations(t)

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
