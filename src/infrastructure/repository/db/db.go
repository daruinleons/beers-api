package db

import (
	"database/sql"
	"fmt"
	"github.com/dleonsal/beers-api/src/configs"
)

const (
	queryCreateBeerTable = `CREATE TABLE IF NOT EXISTS beer (
			id bigint(20) NOT NULL,
			name varchar(45) COLLATE utf8_spanish2_ci DEFAULT NULL,
			brewery varchar(45) COLLATE utf8_spanish2_ci DEFAULT NULL,
			country varchar(45) COLLATE utf8_spanish2_ci NOT NULL,
			price decimal(10,2) DEFAULT NULL,
			currency varchar(32) COLLATE utf8_spanish2_ci NOT NULL,
			PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_spanish2_ci;`
)

func NewMySqlDB(config *configs.DBConfig) *sql.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		config.Username,
		config.Password,
		config.Host,
		config.DBName,
	)

	client, err := sql.Open(config.DriverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	if _, err := client.Exec(queryCreateBeerTable); err != nil {
		panic(err)
	}

	return client
}
