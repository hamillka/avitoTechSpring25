package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechSpring25/internal/handlers"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/middlewares"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/mocks"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func withContextWithRole(role string, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), middlewares.Key("props"), jwt.MapClaims{"role": role})
	return r.WithContext(ctx)
}

func TestAddProductToReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockProductService(ctrl)
	logger := zaptest.NewLogger(t).Sugar()
	handler := handlers.NewProductHandler(service, logger)

	reqBody := dto.AddProductRequestDto{
		Type:  dto.ProductTypeClothes,
		PVZId: "pvz1",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(bodyBytes))
	req = withContextWithRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()

	mockProduct := models.Product{
		Id:          "p1",
		Type:        reqBody.Type,
		ReceptionId: "rec1",
		DateTime:    time.Now().String(),
	}
	service.EXPECT().AddProductToReception(reqBody.Type, reqBody.PVZId).Return(mockProduct, nil)

	handler.AddProductToReception(w, req)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
}

func TestAddProductToReception_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := handlers.NewProductHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/products", nil)
	req = withContextWithRole("moderator", req)
	w := httptest.NewRecorder()

	handler.AddProductToReception(w, req)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestAddProductToReception_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := handlers.NewProductHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader("{invalid json"))
	req = withContextWithRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()

	handler.AddProductToReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAddProductToReception_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := handlers.NewProductHandler(nil, zaptest.NewLogger(t).Sugar())
	reqBody := dto.AddProductRequestDto{
		Type:  "food",
		PVZId: "pvz1",
	}
	data, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(data))
	req = withContextWithRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()

	handler.AddProductToReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAddProductToReception_ErrPVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockProductService(ctrl)
	handler := handlers.NewProductHandler(service, zaptest.NewLogger(t).Sugar())

	reqBody := dto.AddProductRequestDto{
		Type:  dto.ProductTypeClothes,
		PVZId: "pvz404",
	}
	data, _ := json.Marshal(reqBody)

	service.EXPECT().AddProductToReception(reqBody.Type, reqBody.PVZId).Return(models.Product{}, dto.ErrPVZNotFound)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(data))
	req = withContextWithRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()

	handler.AddProductToReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAddProductToReception_ErrNoActiveReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockProductService(ctrl)
	handler := handlers.NewProductHandler(service, zaptest.NewLogger(t).Sugar())

	reqBody := dto.AddProductRequestDto{
		Type:  dto.ProductTypeShoes,
		PVZId: "pvz1",
	}
	data, _ := json.Marshal(reqBody)

	service.EXPECT().AddProductToReception(reqBody.Type, reqBody.PVZId).Return(models.Product{}, dto.ErrNoActiveReception)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(data))
	req = withContextWithRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()

	handler.AddProductToReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAddProductToReception_UnknownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mocks.NewMockProductService(ctrl)
	handler := handlers.NewProductHandler(service, zaptest.NewLogger(t).Sugar())

	reqBody := dto.AddProductRequestDto{
		Type:  dto.ProductTypeShoes,
		PVZId: "pvz1",
	}
	data, _ := json.Marshal(reqBody)

	service.EXPECT().AddProductToReception(reqBody.Type, reqBody.PVZId).Return(models.Product{}, errors.New("unexpected"))

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(data))
	req = withContextWithRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()

	handler.AddProductToReception(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
