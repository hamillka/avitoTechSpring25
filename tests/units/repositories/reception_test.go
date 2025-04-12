package repositories_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hamillka/avitoTechSpring25/internal/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestReceptionRepository_GetLastReception_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)
	timeNow := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1 ORDER BY date_time DESC LIMIT 1`)).
		WithArgs("pvz123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow("rec1", timeNow, "pvz123", "in_progress"))

	r, err := repo.GetLastReception("pvz123")
	assert.NoError(t, err)
	assert.Equal(t, "rec1", r.Id)
	assert.Equal(t, "pvz123", r.PVZId)
	assert.Equal(t, "in_progress", r.Status)
}

func TestReceptionRepository_GetLastReception_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1 ORDER BY date_time DESC LIMIT 1`)).
		WithArgs("pvz123").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetLastReception("pvz123")
	assert.Error(t, err)
}

func TestReceptionRepository_CreateReception_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)
	timeNow := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO receptions (pvz_id) VALUES ($1) RETURNING id, date_time, pvz_id, status`)).
		WithArgs("pvz1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow("rec1", timeNow, "pvz1", "in_progress"))

	r, err := repo.CreateReception("pvz1")
	assert.NoError(t, err)
	assert.Equal(t, "rec1", r.Id)
	assert.Equal(t, "pvz1", r.PVZId)
	assert.Equal(t, "in_progress", r.Status)
}

func TestReceptionRepository_CreateReception_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO receptions (pvz_id) VALUES ($1) RETURNING id, date_time, pvz_id, status`)).
		WithArgs("pvz1").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.CreateReception("pvz1")
	assert.Error(t, err)
}

func TestReceptionRepository_UpdateReceptionStatus_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)
	timeNow := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE receptions SET status = $1 WHERE id = $2 RETURNING id, date_time, pvz_id, status`)).
		WithArgs("close", "rec1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow("rec1", timeNow, "pvz1", "close"))

	r, err := repo.UpdateReceptionStatus("rec1", "close")
	assert.NoError(t, err)
	assert.Equal(t, "close", r.Status)
}

func TestReceptionRepository_UpdateReceptionStatus_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE receptions SET status = $1 WHERE id = $2 RETURNING id, date_time, pvz_id, status`)).
		WithArgs("close", "rec1").
		WillReturnError(sql.ErrTxDone)

	_, err := repo.UpdateReceptionStatus("rec1", "close")
	assert.Error(t, err)
}

func TestReceptionRepository_GetReceptionsByPVZId_NoFilter(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)
	timeNow := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1`)).
		WithArgs("pvz123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow("r1", timeNow, "pvz123", "in_progress"))

	rs, err := repo.GetReceptionsByPVZId("pvz123", nil, nil)
	assert.NoError(t, err)
	assert.Len(t, rs, 1)
	assert.Equal(t, "r1", rs[0].Id)
}

func TestReceptionRepository_GetReceptionsByPVZId_WithFilter(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT r.id, r.date_time, r.pvz_id, r.status 
	FROM receptions r
	WHERE r.pvz_id = $1 
	AND EXISTS (
		SELECT 1 FROM products p 
		WHERE p.reception_id = r.id 
		AND p.date_time BETWEEN $2 AND $3
	)
`)).
		WithArgs("pvz123", start, end).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).
			AddRow("r1", start, "pvz123", "close"))

	rs, err := repo.GetReceptionsByPVZId("pvz123", &start, &end)
	assert.NoError(t, err)
	assert.Len(t, rs, 1)
	assert.Equal(t, "r1", rs[0].Id)
	assert.Equal(t, "close", rs[0].Status)
}

func TestReceptionRepository_GetReceptionsByPVZId_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewReceptionRepository(sqlxDB)
	start := time.Now()
	end := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT r.id, r.date_time, r.pvz_id, r.status 
	FROM receptions r
	WHERE r.pvz_id = $1 
	AND EXISTS (
		SELECT 1 FROM products p 
		WHERE p.reception_id = r.id 
		AND p.date_time BETWEEN $2 AND $3
	)
`)).
		WithArgs("pvz123", start, end).
		WillReturnError(sql.ErrConnDone)

	_, err := repo.GetReceptionsByPVZId("pvz123", &start, &end)
	assert.Error(t, err)
}
