package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/mocks"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestUserHandler_Login_BadJSON(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{bad json"))
	w := httptest.NewRecorder()
	h.Login(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login_InvalidEmail(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	body := dto.UserLoginRequestDto{Email: "invalid", Password: "pass"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.Login(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login_AuthFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserService(ctrl)
	h := NewUserHandler(mockService, zaptest.NewLogger(t).Sugar())

	mockService.EXPECT().UserLogin("test@mail.com", "pass").Return(models.User{}, errors.New("unauthorized"))
	body := dto.UserLoginRequestDto{Email: "test@mail.com", Password: "pass"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.Login(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserService(ctrl)
	h := NewUserHandler(mockService, zaptest.NewLogger(t).Sugar())

	mockService.EXPECT().UserLogin("test@mail.com", "pass").Return(models.User{
		Id:    "1",
		Email: "test@mail.com",
		Role:  "employee",
	}, nil)

	body := dto.UserLoginRequestDto{Email: "test@mail.com", Password: "pass"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.Login(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_Register_BadJSON(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("bad json"))
	w := httptest.NewRecorder()
	h.Register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_InvalidEmail(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	body := dto.UserRegisterRequestDto{Email: "bad", Password: "123", Role: "employee"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.Register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_InvalidRole(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	body := dto.UserRegisterRequestDto{Email: "user@mail.com", Password: "123", Role: "admin"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.Register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserService(ctrl)
	h := NewUserHandler(mockService, zaptest.NewLogger(t).Sugar())

	mockService.EXPECT().UserRegister("user@mail.com", "123", "moderator").Return(models.User{}, errors.New("fail"))

	body := dto.UserRegisterRequestDto{Email: "user@mail.com", Password: "123", Role: "moderator"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.Register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserService(ctrl)
	h := NewUserHandler(mockService, zaptest.NewLogger(t).Sugar())

	mockService.EXPECT().UserRegister("user@mail.com", "123", "moderator").Return(models.User{
		Id:    "u1",
		Email: "user@mail.com",
		Role:  "moderator",
	}, nil)

	body := dto.UserRegisterRequestDto{Email: "user@mail.com", Password: "123", Role: "moderator"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.Register(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserHandler_DummyLogin_InvalidJSON(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	h.DummyLogin(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_DummyLogin_InvalidRole(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	body := dto.DummyLoginRequestDto{Role: "hacker"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.DummyLogin(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_DummyLogin_Success(t *testing.T) {
	h := NewUserHandler(nil, zaptest.NewLogger(t).Sugar())
	body := dto.DummyLoginRequestDto{Role: "employee"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(data))
	w := httptest.NewRecorder()
	h.DummyLogin(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
