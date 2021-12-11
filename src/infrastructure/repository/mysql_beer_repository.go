package repository

import (
	"database/sql"
	genericerrors "errors"
	"fmt"

	"github.com/dleonsal/beers-api/src/core/domain/entities"
	"github.com/dleonsal/beers-api/src/errors"
	"github.com/dleonsal/beers-api/src/infrastructure/logger"
	"github.com/go-sql-driver/mysql"
)

const (
	queryListBeers  = "SELECT id, name, brewery, country, price, currency FROM beer;"
	queryGetBeer    = "SELECT id, name, brewery, country, price, currency FROM beer WHERE id =?"
	queryInsertBeer = "INSERT INTO beer(id, name, brewery, country, price, currency) VALUES(?, ?, ?, ?, ?, ?);"
)

type mySqlBeerRepository struct {
	db *sql.DB
}

func NewMySqlBeerRepository(db *sql.DB) *mySqlBeerRepository {
	return &mySqlBeerRepository{
		db: db,
	}
}

func (r *mySqlBeerRepository) List() ([]entities.Beer, *errors.RestError) {
	stmt, err := r.db.Prepare(queryListBeers)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to prepare statement: %s", err))
		return nil, errors.NewInternalServerError("error trying to get beers from database")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to execute query: %s", err))
		return nil, errors.NewInternalServerError("error trying to get beers from database")
	}
	defer rows.Close()

	beers := make([]entities.Beer, 0)
	for rows.Next() {
		var beer entities.Beer

		if err := rows.Scan(&beer.Id, &beer.Name, &beer.Brewery, &beer.Country, &beer.Price, &beer.Currency); err != nil {
			logger.Log.Error(fmt.Sprintf("error trying to scan rows: %s", err))
			return nil, errors.NewInternalServerError("error trying to get beers from database")
		}

		beers = append(beers, beer)
	}

	return beers, nil
}

func (r *mySqlBeerRepository) GetByID(beerID int64) (*entities.Beer, *errors.RestError) {
	stmt, err := r.db.Prepare(queryGetBeer)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to prepare statement: %s", err))
		return nil, errors.NewInternalServerError("error trying to get beer from database")
	}
	defer stmt.Close()

	bear := entities.Beer{}
	result := stmt.QueryRow(beerID)
	if getErr := result.Scan(&bear.Id, &bear.Name, &bear.Brewery, &bear.Country, &bear.Price, &bear.Currency); getErr != nil {
		if genericerrors.Is(getErr, sql.ErrNoRows) {
			return nil, errors.NewNotFoundError("beer not found")
		}

		logger.Log.Error(fmt.Sprintf("error trying to execute query: %s", getErr))
		return nil, errors.NewInternalServerError("error trying to get beer from database")
	}

	return &bear, nil
}

func (r *mySqlBeerRepository) Save(beer entities.Beer) *errors.RestError {
	stmt, err := r.db.Prepare(queryInsertBeer)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("error trying to prepare statement: %s", err))
		return errors.NewInternalServerError("error trying to save beer in database")
	}
	defer stmt.Close()

	_, saveErr := stmt.Exec(beer.Id, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency)
	if saveErr != nil {
		logger.Log.Error(fmt.Sprintf("error trying to execute query: %s", saveErr))
		driveErr, ok := saveErr.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError("error trying to save beer in database")
		}

		if driveErr.Number == 1062 {
			return errors.NewConflictError(
				fmt.Sprintf("beer id %d already exists", beer.Id))
		}

		return errors.NewInternalServerError("error trying to save beer in database")
	}

	return nil
}
