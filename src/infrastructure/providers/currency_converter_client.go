package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dleonsal/beers-api/src/errors"
	"github.com/dleonsal/beers-api/src/infrastructure/logger"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CurrencyConverterRestClient struct {
	httpClient     HTTPClient
	baseURL        string
	requestTimeout time.Duration
	xAPIKey        string
}

func NewCurrencyConverterRestClient(httpClient HTTPClient, baseURL string, requestTimeout time.Duration, xAPIkey string) *CurrencyConverterRestClient {
	return &CurrencyConverterRestClient{
		httpClient:     httpClient,
		baseURL:        baseURL,
		requestTimeout: requestTimeout,
		xAPIKey:        xAPIkey,
	}
}

func (c *CurrencyConverterRestClient) ConvertValueToNewCurrency(oldCurrency, newCurrency string, value float64) (float64, *errors.RestError) {
	url := fmt.Sprintf("%s/exchange", c.baseURL)
	reqCtx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to create request: %s", err))
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")

	}

	req.Header.Add("x-rapidapi-key", c.xAPIKey)
	q := req.URL.Query()
	q.Add("from", oldCurrency)
	q.Add("to", newCurrency)
	req.URL.RawQuery = q.Encode()

	response, err := c.httpClient.Do(req)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to execute request: %s", err))
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Log.Error(fmt.Sprintf("error trying to convert from one currency to another, response code %d", response.StatusCode))
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")
	}

	var result float64
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to decode response body: %s", err))
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")
	}

	return result * value, nil
}
