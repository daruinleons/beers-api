package configs

import (
	"github.com/dleonsal/beers-api/src/configs/environment"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port                              string                            `yaml:"Port"`
	DBConfig                          DBConfig                          `yaml:"DBConfig"`
	CurrencyConverterRestClientConfig CurrencyConverterRestClientConfig `yaml:"CurrencyConverterRestClientConfig"`
	HTTPClientTimeoutMilliseconds     int                               `yaml:"HTTPClientTimeoutMilliseconds"`
}

type DBConfig struct {
	Username        string `yaml:"UserName"`
	Password        string `yaml:"Password"`
	Host            string `yaml:"Host"`
	DriverName      string `yaml:"DriverName"`
	DBName          string `yaml:"DBName"`
	ConnMaxLifetime int64  `yaml:"ConnMaxLifetime"`
	MaxIdleConns    int    `yaml:"MaxIdleConns"`
	MaxOpenConns    int    `yaml:"MaxOpenConns"`
}

type CurrencyConverterRestClientConfig struct {
	BaseURL                    string `yaml:"BaseURL"`
	RequestTimeoutMilliseconds int    `yaml:"RequestTimeoutMilliseconds"`
	XAPIKey                    string `yaml:"XAPIKey"`
}

func NewConfig() *Config {
	config := new(Config)
	applyConfigFromString(config, []byte(environment.Test))

	return config
}

func applyConfigFromString(config *Config, configBytes []byte) {
	err := yaml.Unmarshal(configBytes, config)

	if err != nil {
		panic(err)
	}
}
