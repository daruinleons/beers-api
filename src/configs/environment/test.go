package environment

const Test = `
Port: 8080
DBConfig:
  UserName: root
  Password: 123456
  Host: mysql-db
  DriverName: mysql
  DBName: BEERSDB
CurrencyConverterRestClientConfig:
  BaseURL: https://currency-exchange.p.rapidapi.com
  RequestTimeoutMilliseconds: 5000
  XAPIKeyEnv: CURRENCY_CONVERTER_X_API_KEY
HTTPClientTimeoutMilliseconds: 5100
`
