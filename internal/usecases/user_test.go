package usecases

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/hamillka/avitoTechSpring25/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(repo)

	repo.EXPECT().UserLogin("test@example.com", "password").Return(models.User{}, errors.New("not found"))

	repo.EXPECT().UserRegister("test@example.com", gomock.Any(), "user").
		Return(models.User{Id: "u1", Email: "test@example.com", Role: "user"}, nil)

	user, err := service.UserRegister("test@example.com", "password", "user")

	require.NoError(t, err)
	assert.Equal(t, "u1", user.Id)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUserRegister_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(repo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	repo.EXPECT().UserLogin("test@example.com", "password").
		Return(models.User{Id: "u1", Email: "test@example.com", Password: string(hashed)}, nil)

	_, err := service.UserRegister("test@example.com", "password", "user")

	assert.ErrorIs(t, err, dto.ErrUserAlreadyExists)
}

func TestUserLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(repo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	repo.EXPECT().UserLogin("test@example.com", "password").
		Return(models.User{Id: "u1", Email: "test@example.com", Password: string(hashed)}, nil)

	user, err := service.UserLogin("test@example.com", "password")

	require.NoError(t, err)
	assert.Equal(t, "u1", user.Id)
}

func TestUserLogin_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	service := NewUserService(repo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)

	repo.EXPECT().UserLogin("test@example.com", "wrongpass").
		Return(models.User{Id: "u1", Email: "test@example.com", Password: string(hashed)}, nil)

	_, err := service.UserLogin("test@example.com", "wrongpass")

	assert.Error(t, err)
}
