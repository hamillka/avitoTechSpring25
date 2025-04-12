//go:generate mockgen -source=pvz.go -destination=./mocks/mock_pvz.go -package=mocks
package usecases

import (
	"time"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
)

type PVZRepository interface {
	CreatePVZ(city string) (models.PVZ, error)
	GetPVZById(pvzId string) (models.PVZ, error)
	GetPVZsWithPagination(offset, limit int) ([]models.PVZ, error)
}

type PVZService struct {
	pvzRepo  PVZRepository
	recRepo  ReceptionRepository
	prodRepo ProductRepository
}

func NewPVZService(pvzRepo PVZRepository, recRepo ReceptionRepository, prodRepo ProductRepository) *PVZService {
	return &PVZService{
		pvzRepo:  pvzRepo,
		recRepo:  recRepo,
		prodRepo: prodRepo,
	}
}

func (pvzs *PVZService) CreatePVZ(city string) (models.PVZ, error) {
	/*
		создать пвз
		вернуть пвз
	*/
	pvz, err := pvzs.pvzRepo.CreatePVZ(city)
	if err != nil {
		return models.PVZ{}, err
	}

	return pvz, nil
}

func (pvzs *PVZService) GetPVZWithPagination(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	/*
		Достать все ПВЗ
		Достать все приемки для этих ПВЗ
		Достать все продукты для этих приемок
		Сгруппировать все по пвз
		Вернуть все ПВЗ с приемками и продуктами
	*/
	offset := (page - 1) * limit

	allPVZs, err := pvzs.pvzRepo.GetPVZsWithPagination(offset, limit)
	if err != nil {
		return nil, err
	}

	result := make([]models.PVZWithReceptions, 0, len(allPVZs))

	dateFiltersApplied := startDate != nil || endDate != nil

	for _, pvz := range allPVZs {
		receptions, err := pvzs.recRepo.GetReceptionsByPVZId(pvz.Id, startDate, endDate)
		if err != nil {
			return nil, err
		}

		if dateFiltersApplied && len(receptions) == 0 {
			continue
		}

		receptionsWithProducts := make([]models.ReceptionWithProducts, 0, len(receptions))

		for _, reception := range receptions {
			products, err := pvzs.prodRepo.GetProductsByReceptionId(reception.Id, startDate, endDate)
			if err != nil {
				return nil, err
			}

			receptionsWithProducts = append(receptionsWithProducts, models.ReceptionWithProducts{
				Reception: reception,
				Products:  products,
			})
		}

		result = append(result, models.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: receptionsWithProducts,
		})
	}

	return result, nil
}

func (pvzs *PVZService) CloseLastReception(pvzId string) (models.Reception, error) {
	/*
		Проверить существует ли пвз. Если нет, то ErrPVZNotFound
		Проверить существует ли приемка у пвз. Если нет, то ErrNoActiveReception
		Закрыть приемку
		Вернуть приемку
	*/
	_, err := pvzs.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	lastReception, err := pvzs.recRepo.GetLastReception(pvzId)
	if err != nil || lastReception.Status == "close" {
		return models.Reception{}, dto.ErrNoActiveReception
	}

	updRec, err := pvzs.recRepo.UpdateReceptionStatus(lastReception.Id, "close")
	if err != nil {
		return models.Reception{}, err
	}

	return updRec, nil
}

func (pvzs *PVZService) DeleteLastProduct(pvzId string) error {
	/*
		Проверить существует ли пвз. Если нет, то ErrPVZNotFound
		Проверить существует ли приемка у пвз и открыта ли она. Если нет, то ErrNoActiveReception
		Проверить существуют ли продукты у приемки. Если нет, то ErrNoProductsInReception
		Удалить продукт из приемки
	*/

	_, err := pvzs.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return err
	}

	lastReception, err := pvzs.recRepo.GetLastReception(pvzId)
	if err != nil || lastReception.Status == "close" {
		return dto.ErrNoActiveReception
	}

	product, err := pvzs.prodRepo.GetLastProduct(lastReception.Id)
	if err != nil {
		return dto.ErrNoProductsInReception
	}

	err = pvzs.prodRepo.DeleteProduct(product.Id)
	if err != nil {
		return err
	}

	return nil
}
