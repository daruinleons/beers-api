package services

import (
	"github.com/dleonsal/beers-api/src/core/domain/entities"
	"github.com/dleonsal/beers-api/src/errors"
	"strings"
)

type BeerRepository interface {
	List() ([]entities.Beer, *errors.RestError)
	GetByID(beerID int64) (*entities.Beer, *errors.RestError)
	Save(beer entities.Beer) *errors.RestError
}

type CurrencyConverterClient interface {
	ConvertValueToNewCurrency(oldCurrency, newCurrency string, value float64) (float64, *errors.RestError)
}

type beerService struct {
	beerRepository          BeerRepository
	currencyConverterClient CurrencyConverterClient
}

func NewBeerService(beerRepository BeerRepository, currencyConverterClient CurrencyConverterClient) *beerService {
	return &beerService{
		beerRepository:          beerRepository,
		currencyConverterClient: currencyConverterClient,
	}
}

func (s *beerService) ListBeers() ([]entities.Beer, *errors.RestError) {
	beers, err := s.beerRepository.List()
	if err != nil {
		return nil, err
	}

	return beers, nil
}

func (s *beerService) GetBeerByID(beerID int64) (*entities.Beer, *errors.RestError) {
	beer, err := s.beerRepository.GetByID(beerID)
	if err != nil {
		return nil, err
	}

	return beer, nil
}

func (s *beerService) GetBoxPrice(beerID int64, newCurrency string, quantity uint64) (float64, *errors.RestError) {
	if len(strings.TrimSpace(newCurrency)) == 0 {
		return 0, errors.NewBadRequestError("currency must not be empty")
	}

	if quantity == uint64(0) {
		quantity = 6
	}

	beer, err := s.GetBeerByID(beerID)
	if err != nil {
		return 0, err
	}

	if beer.Currency == newCurrency {
		totalPrice := beer.Price * float64(quantity)

		return totalPrice, nil
	}

	newPrice, err := s.currencyConverterClient.ConvertValueToNewCurrency(beer.Currency, newCurrency, beer.Price)
	if err != nil {
		return 0, err
	}

	totalPrice := newPrice * float64(quantity)
	return totalPrice, nil
}

func (s *beerService) CreateBeer(beer entities.Beer) *errors.RestError {
	if err := beer.Validate(); err != nil {
		return err
	}

	if err := s.beerRepository.Save(beer); err != nil {
		return err
	}

	return nil
}
