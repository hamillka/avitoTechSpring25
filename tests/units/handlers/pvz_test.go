package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/hamillka/avitoTechSpring25/internal/handlers"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/middlewares"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/mocks"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func withRole(role string, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), middlewares.Key("props"), jwt.MapClaims{"role": role})
	return r.WithContext(ctx)
}

func TestCreatePVZ_Forbidden(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/pvz", nil)
	req = withRole("employee", req)
	w := httptest.NewRecorder()
	handler.CreatePVZ(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreatePVZ_InvalidJSON(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString("{bad json"))
	req = withRole("moderator", req)
	w := httptest.NewRecorder()
	handler.CreatePVZ(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePVZ_InvalidCity(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	body := dto.CreatePVZRequestDto{City: "Berlin"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(data))
	req = withRole("moderator", req)
	w := httptest.NewRecorder()
	handler.CreatePVZ(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePVZ_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())
	service.EXPECT().CreatePVZ(dto.Moscow).Return(models.PVZ{}, errors.New("db error"))
	body := dto.CreatePVZRequestDto{City: dto.Moscow}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(data))
	req = withRole("moderator", req)
	w := httptest.NewRecorder()
	handler.CreatePVZ(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreatePVZ_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())
	pvz := models.PVZ{Id: "123", City: dto.Kazan, RegistrationDate: time.Now().String()}
	service.EXPECT().CreatePVZ(dto.Kazan).Return(pvz, nil)
	body := dto.CreatePVZRequestDto{City: dto.Kazan}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(data))
	req = withRole("moderator", req)
	w := httptest.NewRecorder()
	handler.CreatePVZ(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetPVZWithPagination_InvalidPage(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodGet, "/pvz?page=abc", nil)
	w := httptest.NewRecorder()
	handler.GetPVZWithPagination(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZWithPagination_InvalidLimit(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodGet, "/pvz?limit=999", nil)
	w := httptest.NewRecorder()
	handler.GetPVZWithPagination(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZWithPagination_InvalidDates(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodGet, "/pvz?startDate=wrong", nil)
	w := httptest.NewRecorder()
	handler.GetPVZWithPagination(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZWithPagination_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())
	service.EXPECT().GetPVZWithPagination(nil, nil, 1, 10).Return(nil, errors.New("fail"))
	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	w := httptest.NewRecorder()
	handler.GetPVZWithPagination(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetPVZWithPagination_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())
	service.EXPECT().GetPVZWithPagination(nil, nil, 1, 10).Return([]models.PVZWithReceptions{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
	w := httptest.NewRecorder()
	handler.GetPVZWithPagination(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCloseLastReception_Forbidden(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/pvz/pvz1/close_last_reception", nil)
	req = withRole("moderator", req)
	req = mux.SetURLVars(req, map[string]string{"pvzId": "pvz1"})
	w := httptest.NewRecorder()
	handler.CloseLastReception(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCloseLastReception_BadRequest_NoPVZId(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/pvz/close_last_reception", nil)
	req = withRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()
	handler.CloseLastReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCloseLastReception_ErrNoActiveReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())

	service.EXPECT().CloseLastReception("pvz1").Return(models.Reception{}, dto.ErrNoActiveReception)

	req := httptest.NewRequest(http.MethodPost, "/pvz/pvz1/close_last_reception", nil)
	req = withRole(dto.RoleEmployee, req)
	req = mux.SetURLVars(req, map[string]string{"pvzId": "pvz1"})
	w := httptest.NewRecorder()

	handler.CloseLastReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCloseLastReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())

	rec := models.Reception{Id: "rec1", PVZId: "pvz1", Status: "close", DateTime: time.Now().String()}
	service.EXPECT().CloseLastReception("pvz1").Return(rec, nil)

	req := httptest.NewRequest(http.MethodPost, "/pvz/pvz1/close_last_reception", nil)
	req = withRole(dto.RoleEmployee, req)
	req = mux.SetURLVars(req, map[string]string{"pvzId": "pvz1"})
	w := httptest.NewRecorder()

	handler.CloseLastReception(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteLastProduct_Forbidden(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/pvz/pvz1/delete_last_product", nil)
	req = withRole("moderator", req)
	req = mux.SetURLVars(req, map[string]string{"pvzId": "pvz1"})
	w := httptest.NewRecorder()
	handler.DeleteLastProduct(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteLastProduct_BadRequest_NoPVZId(t *testing.T) {
	handler := handlers.NewPVZHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/pvz/delete_last_product", nil)
	req = withRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()
	handler.DeleteLastProduct(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteLastProduct_ErrNoProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())

	service.EXPECT().DeleteLastProduct("pvz1").Return(dto.ErrNoProductsInReception)

	req := httptest.NewRequest(http.MethodPost, "/pvz/pvz1/delete_last_product", nil)
	req = withRole(dto.RoleEmployee, req)
	req = mux.SetURLVars(req, map[string]string{"pvzId": "pvz1"})
	w := httptest.NewRecorder()

	handler.DeleteLastProduct(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteLastProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockPVZService(ctrl)
	handler := handlers.NewPVZHandler(service, zaptest.NewLogger(t).Sugar())

	service.EXPECT().DeleteLastProduct("pvz1").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/pvz/pvz1/delete_last_product", nil)
	req = withRole(dto.RoleEmployee, req)
	req = mux.SetURLVars(req, map[string]string{"pvzId": "pvz1"})
	w := httptest.NewRecorder()

	handler.DeleteLastProduct(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
