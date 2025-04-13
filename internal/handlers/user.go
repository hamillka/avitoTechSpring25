//go:generate mockgen -source=user.go -destination=./mocks/mock_user.go -package=mocks
package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/middlewares"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"go.uber.org/zap"
)

type UserService interface {
	UserRegister(email, password, role string) (models.User, error)
	UserLogin(email, password string) (models.User, error)
}

type UserHandler struct {
	service UserService
	logger  *zap.SugaredLogger
}

func NewUserHandler(s UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  logger,
	}
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func createToken(role string) (string, error) {
	payload := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(middlewares.Secret)
	if err != nil {
		return "", err
	}

	return t, nil
}

// Login godoc
//
//	@Summary		Авторизация пользователя
//	@Description	Авторизует пользователя по email и паролю и возвращает JWT токен
//	@ID				user-login
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.UserLoginRequestDto	true	"Данные для авторизации"
//
//	@Success		200	{object}	dto.UserLoginResponseDto	"Успешная авторизация"
//	@Failure		400	{object}	dto.ErrorDto				"Некорректные данные"
//	@Failure		401	{object}	dto.ErrorDto				"Неверные учетные данные"
//	@Failure		500	{object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	@Router			/login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userLoginDto dto.UserLoginRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userLoginDto)
	if err != nil {
		uh.logger.Errorf("failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Некорректные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !validateEmail(userLoginDto.Email) {
		uh.logger.Errorf("invalid email format: %v", userLoginDto.Email)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Неверный формат электронной почты",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	user, err := uh.service.UserLogin(userLoginDto.Email, userLoginDto.Password)
	if err != nil {
		uh.logger.Errorf("failed to login user: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		errorDto := &dto.ErrorDto{
			Message: "Неверные учетные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	t, err := createToken(user.Role)
	if err != nil {
		uh.logger.Errorf("failed to create token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Message: "Внутренняя ошибка сервера при авторизации",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	userResponseDto := &dto.UserLoginResponseDto{
		Token: t,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(userResponseDto)
	if err != nil {
		uh.logger.Errorf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Register godoc
//
//	@Summary		Регистрация пользователя
//	@Description	Регистрирует нового пользователя с указанными email, паролем и ролью
//	@ID				user-register
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.UserRegisterRequestDto	true	"Данные для регистрации"
//
//	@Success		201	{object}	dto.UserRegisterResponseDto	"Пользователь успешно зарегистрирован"
//	@Failure		400	{object}	dto.ErrorDto				"Некорректные данные / Неверный запрос"
//	@Failure		500	{object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	@Router			/register [post]
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userRegisterRequestDto dto.UserRegisterRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&userRegisterRequestDto)
	if err != nil {
		uh.logger.Errorf("failed to decode request body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Некорректные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !validateEmail(userRegisterRequestDto.Email) {
		uh.logger.Errorf("invalid email format: %v", userRegisterRequestDto.Email)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Неверный формат электронной почты",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if userRegisterRequestDto.Role != dto.RoleEmployee && userRegisterRequestDto.Role != dto.RoleModerator {
		uh.logger.Errorf("invalid role: %v", userRegisterRequestDto.Role)

		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Неверный запрос",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	user, err := uh.service.UserRegister(
		userRegisterRequestDto.Email,
		userRegisterRequestDto.Password,
		userRegisterRequestDto.Role,
	)
	if err != nil {
		uh.logger.Errorf("failed to register user: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Неверный запрос",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	userRegisterResponseDto := &dto.UserRegisterResponseDto{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(userRegisterResponseDto)
	if err != nil {
		uh.logger.Errorf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// DummyLogin godoc
//
//	@Summary		Упрощенная авторизация
//	@Description	Создает JWT токен с указанной ролью без проверки учетных данных (для тестирования)
//	@ID				user-dummy-login
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.DummyLoginRequestDto	true	"Роль для токена"
//
//	@Success		200	{object}	dto.UserLoginResponseDto	"Успешная авторизация"
//	@Failure		400	{object}	dto.ErrorDto				"Некорректные данные / Неверный запрос"
//	@Failure		500	{object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	@Router			/dummyLogin [post]
func (uh *UserHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var dummyLoginDto dto.DummyLoginRequestDto

	w.Header().Add("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&dummyLoginDto)
	if err != nil {
		uh.logger.Errorf("failed to decode request body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Некорректные данные",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if dummyLoginDto.Role != "employee" && dummyLoginDto.Role != "moderator" {
		uh.logger.Errorf("invalid role: %v", dummyLoginDto.Role)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Неверный запрос",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// CHECK (code duplication with UserLogin)
	t, err := createToken(dummyLoginDto.Role)
	if err != nil {
		uh.logger.Errorf("failed to create token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Message: "Внутренняя ошибка сервера при авторизации",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	userResponseDto := &dto.UserLoginResponseDto{
		Token: t,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(userResponseDto)
	if err != nil {
		uh.logger.Errorf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
