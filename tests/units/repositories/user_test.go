package repositories_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hamillka/avitoTechSpring25/internal/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_UserRegister_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewUserRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, email, password_hash, role`)).
		WithArgs("test@example.com", "hashedpass", "employee").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).
			AddRow("u123", "test@example.com", "hashedpass", "employee"))

	user, err := repo.UserRegister("test@example.com", "hashedpass", "employee")
	assert.NoError(t, err)
	assert.Equal(t, "u123", user.Id)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "employee", user.Role)
}

func TestUserRepository_UserRegister_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewUserRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, email, password_hash, role`)).
		WithArgs("test@example.com", "pass", "employee").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.UserRegister("test@example.com", "pass", "employee")
	assert.Error(t, err)
}

func TestUserRepository_UserLogin_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewUserRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password_hash, role FROM users WHERE email = $1`)).
		WithArgs("login@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).
			AddRow("uid123", "login@example.com", "hashed", "moderator"))

	user, err := repo.UserLogin("login@example.com", "any")
	assert.NoError(t, err)
	assert.Equal(t, "uid123", user.Id)
	assert.Equal(t, "moderator", user.Role)
}

func TestUserRepository_UserLogin_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewUserRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password_hash, role FROM users WHERE email = $1`)).
		WithArgs("nouser@example.com").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.UserLogin("nouser@example.com", "pass")
	assert.Error(t, err)
}

func TestUserRepository_UserLogin_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewUserRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password_hash, role FROM users WHERE email = $1`)).
		WithArgs("user@example.com").
		WillReturnError(sql.ErrConnDone)

	_, err := repo.UserLogin("user@example.com", "pass")
	assert.Error(t, err)
}
