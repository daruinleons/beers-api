package repository_test

import (
	"database/sql"
	genericerrors "errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dleonsal/beers-api/src/core/domain/entities"
	"github.com/dleonsal/beers-api/src/errors"
	"github.com/dleonsal/beers-api/src/infrastructure/repository"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	queryListBeersTest  = "SELECT id, name, brewery, country, price, currency FROM beer;"
	queryGetBeerTest    = "SELECT id, name, brewery, country, price, currency FROM beer WHERE id =?"
	queryInsertBeerTest = "INSERT INTO beer(id, name, brewery, country, price, currency) VALUES(?, ?, ?, ?, ?, ?);"
)

func Test_List_WhenPrepareStmtFail_ThenReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	prepareErr := genericerrors.New("some error")
	expectedError := errors.NewInternalServerError("error trying to get beers from database")
	mock.ExpectPrepare(queryListBeersTest).WillReturnError(prepareErr)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.List()

	assert.Nil(t, beers)
	assert.Equal(t, expectedError, err)
}

func Test_List_WhenExecuteQueryFail_ThenReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryErr := genericerrors.New("some error")
	expectedError := errors.NewInternalServerError("error trying to get beers from database")
	mock.ExpectPrepare(queryListBeersTest)
	mock.ExpectQuery(queryListBeersTest).WillReturnError(queryErr)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.List()

	assert.Nil(t, beers)
	assert.Equal(t, expectedError, err)
}

func Test_List_WhenScanRowsFail_ThenReturnError(t *testing.T) {
	beer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryRows := mock.NewRows([]string{
		"id", "name", "brewery", "country", "price", "currency",
	}).AddRow(beer.Id, beer.Name, beer.Brewery, beer.Country, "invalid", beer.Currency)
	expectedError := errors.NewInternalServerError("error trying to get beers from database")
	mock.ExpectPrepare(queryListBeersTest)
	mock.ExpectQuery(queryListBeersTest).WillReturnRows(queryRows)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.List()

	assert.Nil(t, beers)
	assert.Equal(t, expectedError, err)
}

func Test_List_WhenQueryIsExecutedSuccessfully_ThenReturnBeers(t *testing.T) {
	beer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryRows := mock.NewRows([]string{
		"id", "name", "brewery", "country", "price", "currency",
	}).AddRow(beer.Id, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency)
	expectedBeer := []entities.Beer{*beer}
	mock.ExpectPrepare(queryListBeersTest)
	mock.ExpectQuery(queryListBeersTest).WillReturnRows(queryRows)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.List()

	assert.Equal(t, expectedBeer, beers)
	assert.Nil(t, err)
}

func Test_GetByID_WhenPrepareStmtFail_ThenReturnError(t *testing.T) {
	id := int64(1)
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	prepareErr := genericerrors.New("some error")
	expectedError := errors.NewInternalServerError("error trying to get beer from database")
	mock.ExpectPrepare(queryListBeersTest).WillReturnError(prepareErr)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.GetByID(id)

	assert.Nil(t, beers)
	assert.Equal(t, expectedError, err)
}

func Test_GetByID_WhenExecuteQueryFail_ThenReturnNotFoundError(t *testing.T) {
	id := int64(1)
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryErr := sql.ErrNoRows
	expectedError := errors.NewNotFoundError("beer not found")
	mock.ExpectPrepare(queryGetBeerTest)
	mock.ExpectQuery(queryGetBeerTest).WillReturnError(queryErr)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.GetByID(id)

	assert.Nil(t, beers)
	assert.Equal(t, expectedError, err)
}

func Test_GetByID_WhenExecuteQueryFail_ThenReturnInternalServerError(t *testing.T) {
	id := int64(1)
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryErr := genericerrors.New("some error")
	expectedError := errors.NewInternalServerError("error trying to get beer from database")
	mock.ExpectPrepare(queryGetBeerTest)
	mock.ExpectQuery(queryGetBeerTest).WillReturnError(queryErr)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.GetByID(id)

	assert.Nil(t, beers)
	assert.Equal(t, expectedError, err)
}

func Test_GetByID_WhenQueryIsExecutedSuccessfully_ThenReturnBeer(t *testing.T) {
	id := int64(1)
	expectedBeer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryRows := mock.NewRows([]string{
		"id", "name", "brewery", "country", "price", "currency",
	}).AddRow(expectedBeer.Id, expectedBeer.Name, expectedBeer.Brewery, expectedBeer.Country, expectedBeer.Price, expectedBeer.Currency)
	mock.ExpectPrepare(queryGetBeerTest)
	mock.ExpectQuery(queryGetBeerTest).WillReturnRows(queryRows)
	repo := repository.NewMySqlBeerRepository(db)

	beers, err := repo.GetByID(id)

	assert.Equal(t, expectedBeer, beers)
	assert.Nil(t, err)
}

func Test_Save_WhenPrepareStmtFail_ThenReturnError(t *testing.T) {
	beer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	prepareErr := genericerrors.New("some error")
	expectedError := errors.NewInternalServerError("error trying to save beer in database")
	mock.ExpectPrepare(queryInsertBeerTest).WillReturnError(prepareErr)
	repo := repository.NewMySqlBeerRepository(db)

	err := repo.Save(*beer)

	assert.Equal(t, expectedError, err)
}

func Test_Save_WhenExecuteQueryFailAndCastingErrorToMySQLErrorFail_ThenReturnInternalServerError(t *testing.T) {
	beer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryErr := genericerrors.New("some error")
	expectedError := errors.NewInternalServerError("error trying to save beer in database")
	mock.ExpectPrepare(queryInsertBeerTest)
	mock.ExpectQuery(queryInsertBeerTest).WillReturnError(queryErr)
	repo := repository.NewMySqlBeerRepository(db)

	err := repo.Save(*beer)

	assert.Equal(t, expectedError, err)
}

func Test_Save_WhenExecuteQueryFail_ThenReturnConflictError(t *testing.T) {
	beer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryErr := &mysql.MySQLError{1062, "duplicate entry"}
	expectedError := errors.NewConflictError("beer id 1 already exists")
	mock.ExpectPrepare(queryInsertBeerTest)
	mock.ExpectExec(queryInsertBeerTest).WithArgs(beer.Id, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency).
		WillReturnError(queryErr)
	repo := repository.NewMySqlBeerRepository(db)

	err := repo.Save(*beer)

	assert.Equal(t, expectedError, err)
}

func Test_Save_WhenExecuteQueryFail_ThenReturnInternalServerError(t *testing.T) {
	beer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	queryErr := &mysql.MySQLError{1064, "syntax error"}
	expectedError := errors.NewInternalServerError("error trying to save beer in database")
	mock.ExpectPrepare(queryInsertBeerTest)
	mock.ExpectExec(queryInsertBeerTest).WithArgs(beer.Id, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency).
		WillReturnError(queryErr)
	repo := repository.NewMySqlBeerRepository(db)

	err := repo.Save(*beer)

	assert.Equal(t, expectedError, err)
}


func Test_Save_WhenQueryIsExecutedSuccessfully_ThenReturnNil(t *testing.T) {
	beer := givenBeer()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mock.ExpectPrepare(queryInsertBeerTest)
	mock.ExpectExec(queryInsertBeerTest).WithArgs(beer.Id, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewMySqlBeerRepository(db)

	err := repo.Save(*beer)

	assert.Nil(t, err)
}

func givenBeer() *entities.Beer {
	return &entities.Beer{
		Id:       1,
		Name:     "Pilsen",
		Brewery:  "Bavaria",
		Country:  "Colombia",
		Price:    2500,
		Currency: "COP",
	}
}
