package entities

import (
	"fmt"
	"strings"

	"github.com/dleonsal/beers-api/src/errors"
)

type Beer struct {
	Id       int64   `json:"Id"`
	Name     string  `json:"Name"`
	Brewery  string  `json:"Brewery"`
	Country  string  `json:"Country"`
	Price    float64 `json:"Price"`
	Currency string  `json:"Currency"`
}

func (b *Beer) Validate() *errors.RestError {
	if b.Id == 0 {
		return errors.NewBadRequestError(
			fmt.Sprintf("invalid Id: %d", b.Id))
	}

	if len(strings.TrimSpace(b.Name)) == 0 {
		return errors.NewBadRequestError(
			fmt.Sprintf("invalid Name: %s", b.Name))
	}

	if len(strings.TrimSpace(b.Brewery)) == 0 {
		return errors.NewBadRequestError(
			fmt.Sprintf("invalid Brewery: %s", b.Brewery))
	}

	if len(strings.TrimSpace(b.Country)) == 0 {
		return errors.NewBadRequestError(
			fmt.Sprintf("invalid Country: %s", b.Country))
	}

	if b.Price == 0 {
		return errors.NewBadRequestError(
			fmt.Sprintf("invalid Price: %f", b.Price))
	}

	if len(strings.TrimSpace(b.Currency)) == 0 {
		return errors.NewBadRequestError(
			fmt.Sprintf("invalid Currency: %s", b.Currency))
	}

	return nil
}
