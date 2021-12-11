package configs_test

import (
	"github.com/dleonsal/beers-api/src/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewConfig_ThenReturnConfig(t *testing.T) {
	expectedConfig := configs.Config{
		Port: "8080",
		DBConfig: configs.DBConfig{
			UserName:   "root",
			Password:   "123456",
			Host:       "mysql-db",
			DriverName: "mysql",
			DBName:     "BEERSDB",
		},
		CurrencyConverterRestClientConfig: configs.CurrencyConverterRestClientConfig{
			BaseURL:                    "https://currency-exchange.p.rapidapi.com",
			RequestTimeoutMilliseconds: 5000,
			XAPIKeyEnv:                 "CURRENCY_CONVERTER_X_API_KEY",
		},
		HTTPClientTimeoutMilliseconds: 5100,
	}

	config := configs.NewConfig()

	assert.Equal(t, expectedConfig, *config)
}
