package providers_test

import (
	genericerros "errors"
	"net/http"
	"testing"
	"time"

	"github.com/dleonsal/beers-api/src/errors"
	"github.com/dleonsal/beers-api/src/infrastructure/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	baseURL        = "/test"
	requestTimeout = time.Duration(1000) * time.Millisecond
	xAPIkey        = "123456"
	oldCurrency    = "COP"
	newCurrency    = "USD"
	value          = 10
)

func Test_ConvertValueToNewCurrency_WhenCreateRequestFail_ThenReturnError(t *testing.T) {
	expectedError := errors.NewInternalServerError("error trying to convert from one currency to another")
	client := providers.NewCurrencyConverterRestClient(nil, "%invalid url%", requestTimeout, xAPIkey)

	totalPrice, err := client.ConvertValueToNewCurrency(oldCurrency, newCurrency, value)

	assert.Equal(t, totalPrice, 0.0)
	assert.Equal(t, expectedError, err)
}

func Test_ConvertValueToNewCurrency_WhenDoRequestFail_ThenReturnError(t *testing.T) {
	error := genericerros.New("some error")
	expectedError := errors.NewInternalServerError("error trying to convert from one currency to another")
	mockHTTPClient := new(providers.MockHTTPClient)
	mockHTTPClient.On("Do", mock.Anything).Return(nil, error)
	client := providers.NewCurrencyConverterRestClient(mockHTTPClient, baseURL, requestTimeout, xAPIkey)

	totalPrice, err := client.ConvertValueToNewCurrency(oldCurrency, newCurrency, value)

	assert.Equal(t, totalPrice, 0.0)
	assert.Equal(t, expectedError, err)
}

func Test_ConvertValueToNewCurrency_WhenResponseWithCodeDifferentTo200_ThenReturnError(t *testing.T) {
	expectedResponse := `{"message":"some bad request","error":"bad_request","status":400}`
	expectedError := errors.NewInternalServerError("error trying to convert from one currency to another")
	server := providers.NewMockServerConfig(
		http.StatusBadRequest,
		"application/json",
		nil,
		expectedResponse,
	).CreateMockServer()
	defer server.Close()

	client := providers.NewCurrencyConverterRestClient(&http.Client{}, server.URL, requestTimeout, xAPIkey)

	totalPrice, err := client.ConvertValueToNewCurrency(oldCurrency, newCurrency, value)

	assert.Equal(t, totalPrice, 0.0)
	assert.Equal(t, expectedError, err)
}

func Test_ConvertValueToNewCurrency_WhenDecodeResponseFail_ThenReturnError(t *testing.T) {
	expectedResponse := `{,}`
	expectedError := errors.NewInternalServerError("error trying to convert from one currency to another")
	server := providers.NewMockServerConfig(
		http.StatusOK,
		"application/json",
		nil,
		expectedResponse,
	).CreateMockServer()
	defer server.Close()

	client := providers.NewCurrencyConverterRestClient(&http.Client{}, server.URL, requestTimeout, xAPIkey)

	totalPrice, err := client.ConvertValueToNewCurrency(oldCurrency, newCurrency, value)

	assert.Equal(t, totalPrice, 0.0)
	assert.Equal(t, expectedError, err)
}

func Test_ConvertValueToNewCurrency_WhenProcessIsExecutedSuccessfully_ThenReturnResult(t *testing.T) {
	expectedResponse := `10.0`
	expectedTotalPrice := 100.0
	server := providers.NewMockServerConfig(
		http.StatusOK,
		"application/json",
		nil,
		expectedResponse,
	).CreateMockServer()
	defer server.Close()

	client := providers.NewCurrencyConverterRestClient(&http.Client{}, server.URL, requestTimeout, xAPIkey)

	totalPrice, err := client.ConvertValueToNewCurrency(oldCurrency, newCurrency, value)

	assert.Equal(t, expectedTotalPrice, totalPrice)
	assert.Nil(t, err)
}
