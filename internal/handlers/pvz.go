//go:generate mockgen -source=pvz.go -destination=./mocks/mock_pvz.go -package=mocks
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/middlewares"
	"github.com/hamillka/avitoTechSpring25/internal/metrics"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"go.uber.org/zap"
)

type PVZService interface {
	CreatePVZ(city string) (models.PVZ, error)
	GetPVZWithPagination(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
	CloseLastReception(pvzId string) (models.Reception, error)
	DeleteLastProduct(pvzId string) error
}

type PVZHandler struct {
	service PVZService
	logger  *zap.SugaredLogger
}

func NewPVZHandler(s PVZService, logger *zap.SugaredLogger) *PVZHandler {
	return &PVZHandler{
		service: s,
		logger:  logger,
	}
}

// CreatePVZ godoc
//
//	@Summary		Завести ПВЗ
//	@Description	Создает новый пункт выдачи заказов в указанном городе
//	@ID				create-pvz
//	@Tags			pvz
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.CreatePVZRequestDto	true	"Информация о создаваемом ПВЗ"
//
//	@Success		201	{object}	dto.CreatePVZResponseDto	"ПВЗ успешно создан"
//	@Failure		400	{object}	dto.ErrorDto				"Некорректные данные / Неверный запрос"
//	@Failure		403	{object}	dto.ErrorDto				"Доступ запрещен"
//	@Failure		500	{object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/pvz [post]
func (pvzh *PVZHandler) CreatePVZ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims)
	role := claims["role"].(string)
	if role != dto.RoleModerator {
		pvzh.logger.Errorf("forbidden action : invalid role: %v", role)
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

	var createPVZRequestDto dto.CreatePVZRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&createPVZRequestDto)
	if err != nil {
		pvzh.logger.Errorf("failed to decode request body: %v", err)
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

	if createPVZRequestDto.City != dto.Moscow && createPVZRequestDto.City != dto.SaintP && createPVZRequestDto.City != dto.Kazan {
		pvzh.logger.Errorf("invalid city: %v", createPVZRequestDto.City)
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

	pvz, err := pvzh.service.CreatePVZ(createPVZRequestDto.City)
	if err != nil {
		pvzh.logger.Errorf("failed to create pvz: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Message: "Внутренняя ошибка сервера",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	createPVZResponseDto := dto.CreatePVZResponseDto{
		Id:               pvz.Id,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createPVZResponseDto)
	if err != nil {
		pvzh.logger.Errorf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	metrics.PVZCreated.Inc()
}

// GetPVZWithPagination godoc
//
//	@Summary		Получить список ПВЗ с пагинацией
//	@Description	Возвращает список ПВЗ с их приемками и товарами с возможностью фильтрации по дате и пагинацией
//	@ID				get-pvz-with-pagination
//	@Tags			pvz
//	@Accept			json
//	@Produce		json
//	@Param			startDate	query	string	false	"Начальная дата (RFC3339)"
//	@Param			endDate		query	string	false	"Конечная дата (RFC3339)"
//	@Param			page		query	integer	false	"Номер страницы (по умолчанию 1)"
//	@Param			limit		query	integer	false	"Количество элементов на странице (по умолчанию 10, максимум 30)"
//
//	@Success		200	{array}		dto.PVZWithReceptionsDto	"Список ПВЗ с приемками и товарами"
//	@Failure		400	{object}	dto.ErrorDto				"Невалидные параметры запроса"
//	@Failure		500	{object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/pvz [get]
func (pvzh *PVZHandler) GetPVZWithPagination(w http.ResponseWriter, r *http.Request) {
	startDateStr, _ := GetQueryParam(r, "startDate", "")
	endDateStr, _ := GetQueryParam(r, "endDate", "")

	page, err := GetQueryParam(r, "page", 1)
	if err != nil || page < 1 {
		pvzh.logger.Errorf("error in extracting page from query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Невалидный параметр page",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	limit, err := GetQueryParam(r, "limit", 10)
	if err != nil || limit < 1 || limit > 30 {
		pvzh.logger.Errorf("error in extracting limit from query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Невалидный параметр limit",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var startDate, endDate *time.Time
	var tStart, tEnd time.Time
	if startDateStr != "" {
		tStart, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			pvzh.logger.Errorf("startDate invalid format: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto := &dto.ErrorDto{
				Message: "Неверный формат даты. Используйте формат RFC3339: 2025-04-11T18:57:00+03:00",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		startDate = &tStart
	}

	if endDateStr != "" {
		tEnd, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			pvzh.logger.Errorf("endDate invalid format: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto := &dto.ErrorDto{
				Message: "Неверный формат даты. Используйте формат RFC3339: 2025-04-11T18:57:00+03:00",
			}
			err = json.NewEncoder(w).Encode(errorDto)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		endDate = &tEnd
	}

	if startDate != nil && endDate != nil && !startDate.Before(*endDate) {
		pvzh.logger.Errorf("startDate should be before endDate")
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

	pvzs, err := pvzh.service.GetPVZWithPagination(startDate, endDate, page, limit)
	if err != nil {
		pvzh.logger.Errorf("failed to get pvzs: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorDto := &dto.ErrorDto{
			Message: "Внутренняя ошибка сервера",
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	pvzsWithReceptionsDto := dto.PVZConvertBLtoDto(pvzs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pvzsWithReceptionsDto)
	if err != nil {
		pvzh.logger.Errorf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// CloseLastReception godoc
//
//	@Summary		Закрыть последнюю приемку
//	@Description	Закрывает последнюю активную приемку для указанного ПВЗ
//	@ID				close-last-reception
//	@Tags			pvz
//	@Accept			json
//	@Produce		json
//	@Param			pvzId	path	string	true	"Идентификатор ПВЗ"
//
//	@Success		200	{object}	dto.CloseReceptionResponseDto	"Приемка успешно закрыта"
//	@Failure		400	{object}	dto.ErrorDto					"Некорректные данные / ПВЗ не найден / Приемка уже закрыта"
//	@Failure		403	{object}	dto.ErrorDto					"Доступ запрещен"
//	@Failure		500	{object}	dto.ErrorDto					"Внутренняя ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/pvz/{pvzId}/close_last_reception [post]
func (pvzh *PVZHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims)
	role := claims["role"].(string)
	if role != dto.RoleEmployee {
		pvzh.logger.Errorf("forbidden action : invalid role: %v", role)
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

	pvzId, ok := mux.Vars(r)["pvzId"]
	if !ok {
		pvzh.logger.Errorf("failed to extract pvzId")
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Некорректные данные",
		}
		err := json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	reception, err := pvzh.service.CloseLastReception(pvzId)
	if err != nil {
		pvzh.logger.Errorf("failed to close last reception: %v", err)
		var errorDto *dto.ErrorDto
		if errors.Is(err, dto.ErrPVZNotFound) {
			pvzh.logger.Errorf("pvz not found: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "ПВЗ не найден",
			}
		} else if errors.Is(err, dto.ErrNoActiveReception) {
			pvzh.logger.Errorf("all receptions are closed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "Приемка уже закрыта",
			}
		} else {
			pvzh.logger.Errorf("internal error: %v", err)
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

	closeReceptionResponseDto := dto.CloseReceptionResponseDto{
		Id:       reception.Id,
		DateTime: reception.DateTime,
		PVZId:    reception.PVZId,
		Status:   reception.Status,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(closeReceptionResponseDto)
	if err != nil {
		pvzh.logger.Errorf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// DeleteLastProduct godoc
//
//	@Summary		Удалить последний товар
//	@Description	Удаляет последний добавленный товар из последней активной приемки для указанного ПВЗ
//	@ID				delete-last-product
//	@Tags			pvz
//	@Accept			json
//	@Produce		json
//	@Param			pvzId	path	string	true	"Идентификатор ПВЗ"
//
//	@Success		200	{object}	nil							"Товар успешно удален"
//	@Failure		400	{object}	dto.ErrorDto				"Некорректные данные / ПВЗ не найден / Нет активной приемки / Нет товаров для удаления"
//	@Failure		403	{object}	dto.ErrorDto				"Доступ запрещен"
//	@Failure		500	{object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/pvz/{pvzId}/delete_last_product [post]
func (pvzh *PVZHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims)
	role := claims["role"].(string)
	if role != dto.RoleEmployee {
		pvzh.logger.Errorf("forbidden action : invalid role: %v", role)
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

	pvzId, ok := mux.Vars(r)["pvzId"]
	if !ok {
		pvzh.logger.Errorf("failed to extract pvzId")
		w.WriteHeader(http.StatusBadRequest)
		errorDto := &dto.ErrorDto{
			Message: "Некорректные данные",
		}
		err := json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err := pvzh.service.DeleteLastProduct(pvzId)
	if err != nil {
		pvzh.logger.Errorf("failed to delete last product: %v", err)
		var errorDto *dto.ErrorDto
		if errors.Is(err, dto.ErrPVZNotFound) {
			pvzh.logger.Errorf("pvz not found: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "ПВЗ не найден",
			}
		} else if errors.Is(err, dto.ErrNoActiveReception) {
			pvzh.logger.Errorf("no active receptions: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "Нет активной приемки",
			}
		} else if errors.Is(err, dto.ErrNoProductsInReception) {
			pvzh.logger.Errorf("no products in receptions: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "Нет товаров для удаления",
			}
		}
		err = json.NewEncoder(w).Encode(errorDto)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetQueryParam[T any](r *http.Request, key string, defaultValue T) (T, error) {
	value := r.URL.Query().Get(key)

	if value == "" {
		return defaultValue, nil
	}

	switch any(defaultValue).(type) {
	case string:
		return any(value).(T), nil

	case int:
		num, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue, err
		}
		return any(num).(T), nil

	default:
		return defaultValue, nil
	}
}
