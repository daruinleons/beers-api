package entities_test

import (
	"fmt"
	"github.com/dleonsal/beers-api/src/core/domain/entities"
	"github.com/dleonsal/beers-api/src/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Validate_WhenIdIsInvalid_ThenReturnBadRequestError(t *testing.T) {
	beer := entities.Beer{
		Id: 0,
	}
	expectedError := errors.NewBadRequestError(fmt.Sprintf("invalid Id: %d", beer.Id))

	err := beer.Validate()

	assert.Equal(t, expectedError, err)
}

func Test_Validate_WhenNameIsInvalid_ThenReturnBadRequestError(t *testing.T) {
	beer := entities.Beer{
		Id:   1,
		Name: "",
	}
	expectedError := errors.NewBadRequestError(fmt.Sprintf("invalid Name: %s", beer.Name))

	err := beer.Validate()

	assert.Equal(t, expectedError, err)
}

func Test_Validate_WhenBreweryIsInvalid_ThenReturnBadRequestError(t *testing.T) {
	beer := entities.Beer{
		Id:      1,
		Name:    "Pilsen",
		Brewery: "",
	}
	expectedError := errors.NewBadRequestError(fmt.Sprintf("invalid Brewery: %s", beer.Brewery))

	err := beer.Validate()

	assert.Equal(t, expectedError, err)
}

func Test_Validate_WhenCountryIsInvalid_ThenReturnBadRequestError(t *testing.T) {
	beer := entities.Beer{
		Id:      1,
		Name:    "Pilsen",
		Brewery: "Bavaria",
		Country: "",
	}
	expectedError := errors.NewBadRequestError(fmt.Sprintf("invalid Country: %s", beer.Country))

	err := beer.Validate()

	assert.Equal(t, expectedError, err)
}

func Test_Validate_WhenPriceIsInvalid_ThenReturnBadRequestError(t *testing.T) {
	beer := entities.Beer{
		Id:      1,
		Name:    "Pilsen",
		Brewery: "Bavaria",
		Country: "Colombia",
		Price:   0.0,
	}
	expectedError := errors.NewBadRequestError(fmt.Sprintf("invalid Price: %f", beer.Price))

	err := beer.Validate()

	assert.Equal(t, expectedError, err)
}

func Test_Validate_WhenCurrencyIsInvalid_ThenReturnBadRequestError(t *testing.T) {
	beer := entities.Beer{
		Id:       1,
		Name:     "Pilsen",
		Brewery:  "Bavaria",
		Country:  "Colombia",
		Price:    2500,
		Currency: "",
	}
	expectedError := errors.NewBadRequestError(fmt.Sprintf("invalid Currency: %s", beer.Currency))

	err := beer.Validate()

	assert.Equal(t, expectedError, err)
}


func Test_Validate_WhenBeerIsValid_ThenReturnNil(t *testing.T) {
	beer := entities.Beer{
		Id:       1,
		Name:     "Pilsen",
		Brewery:  "Bavaria",
		Country:  "Colombia",
		Price:    2500,
		Currency: "COP",
	}

	err := beer.Validate()

	assert.Nil(t, err)
}