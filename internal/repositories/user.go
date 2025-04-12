package repositories

import (
	"database/sql"
	"errors"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

const (
	createUser     = "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, email, password_hash, role"
	getUserByEmail = "SELECT id, email, password_hash, role FROM users WHERE email = $1"
)

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) UserRegister(email, password, role string) (models.User, error) {
	var user models.User
	err := ur.db.QueryRow(createUser, email, password, role).
		Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.Role,
		)
	if err != nil {
		return models.User{}, dto.ErrDBInsert
	}

	return user, nil
}

func (ur *UserRepository) UserLogin(email, password string) (models.User, error) {
	var user models.User

	err := ur.db.QueryRow(getUserByEmail, email).
		Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.Role,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, dto.ErrInvalidCredentials
		}
		return models.User{}, dto.ErrDBRead
	}

	return user, nil
}
