//go:generate mockgen -source=reception.go -destination=./mocks/mock_reception.go -package=mocks
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/middlewares"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"go.uber.org/zap"
)

type ReceptionService interface {
	CreateReception(pvzId string) (models.Reception, error)
}

type ReceptionHandler struct {
	service ReceptionService
	logger  *zap.SugaredLogger
}

func NewReceptionHandler(service ReceptionService, logger *zap.SugaredLogger) *ReceptionHandler {
	return &ReceptionHandler{
		service: service,
		logger:  logger,
	}
}

// CreateReception godoc
//
//	@Summary		Создать приемку
//	@Description	Создает новую приемку товаров для указанного ПВЗ
//	@ID				create-reception
//	@Tags			receptions
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.CreateReceptionRequestDto	true	"Информация о создаваемой приемке"
//
//	@Success		201	{object}	dto.CreateReceptionResponseDto	"Приемка успешно создана"
//	@Failure		400	{object}	dto.ErrorDto					"Некорректные данные / ПВЗ не найден / ПВЗ уже имеет незакрытую приемку"
//	@Failure		403	{object}	dto.ErrorDto					"Доступ запрещен"
//	@Failure		500	{object}	dto.ErrorDto					"Внутренняя ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/receptions [post]
func (rh *ReceptionHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims)
	role := claims["role"].(string)
	if role != dto.RoleEmployee {
		rh.logger.Errorf("forbidden action : invalid role: %v", role)
		w.WriteHeader(http.StatusForbidden)
		errorDto := &dto.ErrorDto{
			Message: "Доступ запрещен",
		}
		err := json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var createReceptionRequestDto dto.CreateReceptionRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&createReceptionRequestDto)
	if err != nil {
		rh.logger.Errorf("failed to decode request body: %v", err)
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

	reception, err := rh.service.CreateReception(createReceptionRequestDto.PVZId)
	if err != nil {
		var errorDto *dto.ErrorDto
		if errors.Is(err, dto.ErrPVZNotFound) {
			rh.logger.Errorf("failed to create reception: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "ПВЗ не найден",
			}
		} else if errors.Is(err, dto.ErrPVZAlreadyHasReception) {
			rh.logger.Errorf("failed to create reception: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "ПВЗ уже имеет незакрытую приемку",
			}
		} else {
			rh.logger.Errorf("failed to create reception: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			errorDto = &dto.ErrorDto{
				Message: "Внутренняя ошибка сервера",
			}
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	createReceptionResponseDto := &dto.CreateReceptionResponseDto{
		Id:       reception.Id,
		DateTime: reception.DateTime,
		PVZId:    reception.PVZId,
		Status:   reception.Status,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createReceptionResponseDto)
	if err != nil {
		rh.logger.Errorf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
