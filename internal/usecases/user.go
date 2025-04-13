//go:generate mockgen -source=user.go -destination=./mocks/mock_user.go -package=mocks
package usecases

import (
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	UserRegister(email, password, role string) (models.User, error)
	UserLogin(email, password string) (models.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (us *UserService) UserRegister(email, password, role string) (models.User, error) {
	existingUser, err := us.userRepo.UserLogin(email, password)
	if err == nil && existingUser.Id != "" {
		return models.User{}, dto.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user, err := us.userRepo.UserRegister(email, string(hashedPassword), role)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (us *UserService) UserLogin(email, password string) (models.User, error) {
	user, err := us.userRepo.UserLogin(email, password)
	if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
