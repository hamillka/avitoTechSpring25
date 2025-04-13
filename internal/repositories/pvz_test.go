package repositories

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPVZRepository_CreatePVZ_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPVZRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO pvzs (city) VALUES ($1) RETURNING id, registration_date, city`)).
		WithArgs("Москва").
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow("123", time.Now(), "Москва"))

	pvz, err := repo.CreatePVZ("Москва")
	assert.NoError(t, err)
	assert.Equal(t, "123", pvz.Id)
	assert.Equal(t, "Москва", pvz.City)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPVZRepository_CreatePVZ_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPVZRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO pvzs (city) VALUES ($1) RETURNING id, registration_date, city`)).
		WithArgs("Казань").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.CreatePVZ("Казань")
	assert.Error(t, err)
}

func TestPVZRepository_GetPVZById_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPVZRepository(sqlxDB)
	timeNow := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, registration_date, city FROM pvzs WHERE id = $1`)).
		WithArgs("abc123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow("abc123", timeNow, "Санкт-Петербург"))

	pvz, err := repo.GetPVZById("abc123")
	assert.NoError(t, err)
	assert.Equal(t, "abc123", pvz.Id)
	assert.Equal(t, "Санкт-Петербург", pvz.City)
}

func TestPVZRepository_GetPVZById_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPVZRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, registration_date, city FROM pvzs WHERE id = $1`)).
		WithArgs("notfound").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetPVZById("notfound")
	assert.Error(t, err)
}

func TestPVZRepository_GetPVZsWithPagination_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPVZRepository(sqlxDB)
	timeNow := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, registration_date, city FROM pvzs ORDER BY registration_date DESC LIMIT $1 OFFSET $2`)).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow("id1", timeNow, "Москва").
			AddRow("id2", timeNow, "Казань"))

	pvzs, err := repo.GetPVZsWithPagination(0, 10)
	assert.NoError(t, err)
	assert.Len(t, pvzs, 2)
	assert.Equal(t, "Москва", pvzs[0].City)
	assert.Equal(t, "Казань", pvzs[1].City)
}

func TestPVZRepository_GetPVZsWithPagination_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewPVZRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, registration_date, city FROM pvzs ORDER BY registration_date DESC LIMIT $1 OFFSET $2`)).
		WithArgs(10, 0).
		WillReturnError(sql.ErrConnDone)

	_, err := repo.GetPVZsWithPagination(0, 10)
	assert.Error(t, err)
}
