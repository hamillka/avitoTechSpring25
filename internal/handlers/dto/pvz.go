package dto

import "github.com/hamillka/avitoTechSpring25/internal/models"

const (
	Moscow = "Москва"
	SaintP = "Санкт-Петербург"
	Kazan  = "Казань"
)

// PVZDto model info
// @Description Информация о ПВЗ
type PVZDto struct {
	Id               string `json:"id"`               // Идентификатор
	RegistrationDate string `json:"registrationDate"` // Дата регистрации
	City             string `json:"city"`             // Город
}

// CreatePVZRequestDto model info
// @Description Информация о ПВЗ при его создании
type CreatePVZRequestDto struct {
	City string `json:"city"` // Город
}

// CreatePVZResponseDto model info
// @Description Информация о ПВЗ при его создании
type CreatePVZResponseDto struct {
	Id               string `json:"id"`               // Идентификатор
	RegistrationDate string `json:"registrationDate"` // Дата регистрации
	City             string `json:"city"`             // Город
}

// PVZWithReceptions model info
// @Description Информация о ПВЗ и приемках, связанных с ним
type PVZWithReceptionsDto struct {
	PVZ        PVZDto                     `json:"pvz"`        // Информация о ПВЗ
	Receptions []ReceptionWithProductsDto `json:"receptions"` // Информация о всех приемках на ПВЗ
}

func PVZConvertBLtoDto(pvzs []models.PVZWithReceptions) []PVZWithReceptionsDto {
	result := make([]PVZWithReceptionsDto, 0, len(pvzs))

	for _, pvz := range pvzs {
		pvzDto := PVZDto{
			Id:               pvz.PVZ.Id,
			RegistrationDate: pvz.PVZ.RegistrationDate,
			City:             pvz.PVZ.City,
		}

		receptionsDto := make([]ReceptionWithProductsDto, 0, len(pvz.Receptions))

		for _, reception := range pvz.Receptions {
			receptionDto := ReceptionDto{
				Id:       reception.Reception.Id,
				DateTime: reception.Reception.DateTime,
				PVZId:    reception.Reception.PVZId,
				Status:   reception.Reception.Status,
			}

			productsDto := make([]ProductDto, 0, len(reception.Products))

			for _, product := range reception.Products {
				productDto := ProductDto{
					Id:          product.Id,
					DateTime:    product.DateTime,
					Type:        product.Type,
					ReceptionId: product.ReceptionId,
				}

				productsDto = append(productsDto, productDto)
			}

			receptionWithProductsDto := ReceptionWithProductsDto{
				Reception: receptionDto,
				Products:  productsDto,
			}

			receptionsDto = append(receptionsDto, receptionWithProductsDto)
		}

		pvzWithReceptionsDto := PVZWithReceptionsDto{
			PVZ:        pvzDto,
			Receptions: receptionsDto,
		}

		result = append(result, pvzWithReceptionsDto)
	}

	return result
}
