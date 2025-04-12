//go:generate mockgen -source=product.go -destination=./mocks/mock_product.go -package=mocks
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

type ProductService interface {
	AddProductToReception(productType, pvzId string) (models.Product, error)
}

type ProductHandler struct {
	service ProductService
	logger  *zap.SugaredLogger
}

func NewProductHandler(service ProductService, logger *zap.SugaredLogger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

// AddProductToReception godoc
//
//	@Summary		Добавить товар в приемку
//	@Description	Добавляет новый товар в активную приемку для указанного ПВЗ
//	@ID				add-product-to-reception
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			body	body	dto.AddProductRequestDto	true	"Информация о товаре"
//
//	@Success		201	{object}	dto.AddProductResponseDto	"Товар успешно добавлен"
//	@Failure		400	{object}	dto.ErrorDto				"Некорректные данные / ПВЗ не найден / Нет активной приемки"
//	@Failure		403	{object}	dto.ErrorDto				"Доступ запрещен"
//	@Failure		500	{object}	dto.ErrorDto				"Внутренняя ошибка сервера"
//	@Security		ApiKeyAuth
//	@Router			/products [post]
func (ph *ProductHandler) AddProductToReception(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := middlewares.Key("props")
	claims := ctx.Value(key).(jwt.MapClaims)
	role := claims["role"].(string)
	if role != dto.RoleEmployee {
		ph.logger.Errorf("forbidden action : invalid role: %v", role)
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

	var addProductRequestDto dto.AddProductRequestDto

	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&addProductRequestDto)
	if err != nil {
		ph.logger.Errorf("failed to decode request body %v", err)
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

	if addProductRequestDto.Type != dto.ProductTypeElectronics && addProductRequestDto.Type != dto.ProductTypeClothes && addProductRequestDto.Type != dto.ProductTypeShoes {
		ph.logger.Errorf("invalid product type: %v", addProductRequestDto.Type)
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

	product, err := ph.service.AddProductToReception(addProductRequestDto.Type, addProductRequestDto.PVZId)
	if err != nil {
		ph.logger.Errorf("failed to add product to reception: %v", err)
		var errorDto *dto.ErrorDto
		if errors.Is(err, dto.ErrPVZNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "ПВЗ не найден",
			}
		} else if errors.Is(err, dto.ErrNoActiveReception) {
			w.WriteHeader(http.StatusBadRequest)
			errorDto = &dto.ErrorDto{
				Message: "Нет активной приемки",
			}
		} else {
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

	addProductResponseDto := &dto.AddProductResponseDto{
		Id:          product.Id,
		DateTime:    product.DateTime,
		Type:        product.Type,
		ReceptionId: product.ReceptionId,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(addProductResponseDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
