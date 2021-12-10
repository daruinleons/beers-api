package environment

const Test = `
Port: 8080
DBConfig:
  UserName: root
  Password: 123456
  Host: mysql-db
  DriverName: mysql
  DBName: beers_db
  ConnMaxLifetime: 300000000000 #5 minutes in nanoseconds
  MaxIdleConns: 2
  MaxOpenConns: 4
CurrencyConverterRestClientConfig:
  BaseURL: https://currency-exchange.p.rapidapi.com
  RequestTimeoutMilliseconds: 5000
  XAPIKey: 84f1ba3f7emshc8742a057cd7a09p1df129jsn9f746d7e4b0d
HTTPClientTimeoutMilliseconds: 5100
`
