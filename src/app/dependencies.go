package app

import (
	"net/http"
	"os"
	"time"

	"github.com/dleonsal/beers-api/src/configs"
	"github.com/dleonsal/beers-api/src/core/services"
	"github.com/dleonsal/beers-api/src/infrastructure/handler"
	"github.com/dleonsal/beers-api/src/infrastructure/providers"
	"github.com/dleonsal/beers-api/src/infrastructure/repository"
	"github.com/dleonsal/beers-api/src/infrastructure/repository/db"
)

const (
	mysqlUsersUsername = "mysql_username"
	mysqlUsersPassword = "mysql_password"
	mysqlUsersHost     = "mysql_host"
	mysqlUsersSchema   = "mysql_schema"
)

var (
	username = os.Getenv(mysqlUsersUsername)
	password = os.Getenv(mysqlUsersPassword)
	host     = os.Getenv(mysqlUsersHost)
	schema   = os.Getenv(mysqlUsersSchema)
)

func wireDependencies(config *configs.Config) *handlerContainer {
	client := db.NewMySqlDB(&config.DBConfig)
	beerRepository := repository.NewMySqlBeerRepository(client)
	httpClient := &http.Client{
		Timeout: time.Duration(config.HTTPClientTimeoutMilliseconds) * time.Millisecond,
	}

	currencyConverterClient := providers.NewCurrencyConverterRestClient(
		httpClient,
		config.CurrencyConverterRestClientConfig.BaseURL,
		time.Duration(config.CurrencyConverterRestClientConfig.RequestTimeoutMilliseconds)*time.Millisecond,
		os.Getenv(config.CurrencyConverterRestClientConfig.XAPIKeyEnv))
	beerService := services.NewBeerService(beerRepository, currencyConverterClient)
	beerHandler := handler.NewBeerHandler(beerService)

	return newHandlerContainer(beerHandler)
}
