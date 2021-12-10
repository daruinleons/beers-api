package providers

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dleonsal/beers-api/src/errors"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CurrencyConverterRestClient struct {
	httpClient     HTTPClient
	baseURL        string
	requestTimeout time.Duration
	XAPIKey        string
}

func NewCurrencyConverterRestClient(httpClient HTTPClient, baseURL string, requestTimeout time.Duration, xAPIkey string) *CurrencyConverterRestClient {
	return &CurrencyConverterRestClient{
		httpClient:     httpClient,
		baseURL:        baseURL,
		requestTimeout: requestTimeout,
		XAPIKey:        xAPIkey,
	}
}

func (c *CurrencyConverterRestClient) ConvertValueToNewCurrency(oldCurrency, newCurrency string, value float64) (float64, *errors.RestError) {
	url := fmt.Sprintf("%s/exchange", c.baseURL)
	reqCtx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")

	}

	req.Header.Add("x-rapidapi-key", c.XAPIKey)
	q := req.URL.Query()
	q.Add("from", oldCurrency)
	q.Add("to", newCurrency)
	req.URL.RawQuery = q.Encode()

	response, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")
	}

	var result float64
	err = json.NewDecoder(response.Body).Decode(&result)
	if response.StatusCode != http.StatusOK {
		fmt.Println(err)
		return 0, errors.NewInternalServerError("error trying to convert from one currency to another")
	}

	return result * value, nil
}
