package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/mocks"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestCreateReception_Forbidden(t *testing.T) {
	handler := NewReceptionHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/receptions", nil)
	req = withRole("moderator", req)
	w := httptest.NewRecorder()
	handler.CreateReception(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateReception_InvalidJSON(t *testing.T) {
	handler := NewReceptionHandler(nil, zaptest.NewLogger(t).Sugar())
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBufferString("{bad json"))
	req = withRole(dto.RoleEmployee, req)
	w := httptest.NewRecorder()
	handler.CreateReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateReception_ErrPVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockReceptionService(ctrl)
	handler := NewReceptionHandler(service, zaptest.NewLogger(t).Sugar())

	reqDto := dto.CreateReceptionRequestDto{PVZId: "pvz404"}
	jsonBody, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(jsonBody))
	req = withRole(dto.RoleEmployee, req)

	service.EXPECT().CreateReception(reqDto.PVZId).Return(models.Reception{}, dto.ErrPVZNotFound)

	w := httptest.NewRecorder()
	handler.CreateReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateReception_ErrAlreadyHasReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockReceptionService(ctrl)
	handler := NewReceptionHandler(service, zaptest.NewLogger(t).Sugar())

	reqDto := dto.CreateReceptionRequestDto{PVZId: "pvz1"}
	jsonBody, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(jsonBody))
	req = withRole(dto.RoleEmployee, req)

	service.EXPECT().CreateReception(reqDto.PVZId).Return(models.Reception{}, dto.ErrPVZAlreadyHasReception)

	w := httptest.NewRecorder()
	handler.CreateReception(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateReception_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockReceptionService(ctrl)
	handler := NewReceptionHandler(service, zaptest.NewLogger(t).Sugar())

	reqDto := dto.CreateReceptionRequestDto{PVZId: "pvz1"}
	jsonBody, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(jsonBody))
	req = withRole(dto.RoleEmployee, req)

	service.EXPECT().CreateReception(reqDto.PVZId).Return(models.Reception{}, errors.New("db failure"))

	w := httptest.NewRecorder()
	handler.CreateReception(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mocks.NewMockReceptionService(ctrl)
	handler := NewReceptionHandler(service, zaptest.NewLogger(t).Sugar())

	reception := models.Reception{Id: "rec123", PVZId: "pvz1", Status: "in_progress", DateTime: time.Now().String()}
	reqDto := dto.CreateReceptionRequestDto{PVZId: "pvz1"}
	jsonBody, _ := json.Marshal(reqDto)
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(jsonBody))
	req = withRole(dto.RoleEmployee, req)

	service.EXPECT().CreateReception(reqDto.PVZId).Return(reception, nil)

	w := httptest.NewRecorder()
	handler.CreateReception(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
