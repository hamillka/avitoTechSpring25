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

func TestAddProduct_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	time := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO products (product_type, reception_id) VALUES ($1, $2) RETURNING id, date_time, product_type, reception_id")).
		WithArgs("одежда", "rec1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "product_type", "reception_id"}).
			AddRow("prod1", time, "одежда", "rec1"))

	p, err := repo.AddProduct("одежда", "rec1")
	assert.NoError(t, err)
	assert.Equal(t, "prod1", p.Id)
}

func TestAddProduct_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO products (product_type, reception_id) VALUES ($1, $2) RETURNING id, date_time, product_type, reception_id")).
		WithArgs("одежда", "rec1").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.AddProduct("одежда", "rec1")
	assert.Error(t, err)
}

func TestGetLastProduct_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	time := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1")).
		WithArgs("rec1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "product_type", "reception_id"}).
			AddRow("prod1", time, "одежда", "rec1"))

	p, err := repo.GetLastProduct("rec1")
	assert.NoError(t, err)
	assert.Equal(t, "prod1", p.Id)
}

func TestGetLastProduct_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1")).
		WithArgs("rec1").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetLastProduct("rec1")
	assert.Error(t, err)
}

func TestDeleteProduct_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
		WithArgs("prod1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteProduct("prod1")
	assert.NoError(t, err)
}

func TestDeleteProduct_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
		WithArgs("prod1").
		WillReturnError(sql.ErrConnDone)

	err := repo.DeleteProduct("prod1")
	assert.Error(t, err)
}

func TestGetProductsByReceptionIds_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	time := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, date_time, product_type, reception_id FROM products WHERE reception_id = ANY($1) AND date_time BETWEEN $2 AND $3")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "product_type", "reception_id"}).
			AddRow("prod1", time, "обувь", "rec1").
			AddRow("prod2", time, "одежда", "rec2"))

	products, err := repo.GetProductsByReceptionIds([]string{"rec1", "rec2"}, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "prod1", products[0].Id)
	assert.Equal(t, "prod2", products[1].Id)
}

func TestGetProductsByReceptionIds_ErrorInQuery(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewProductRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, date_time, product_type, reception_id FROM products WHERE reception_id = ANY($1) AND date_time BETWEEN $2 AND $3")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	products, err := repo.GetProductsByReceptionIds([]string{"rec1", "rec2"}, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, products, 0)
}
