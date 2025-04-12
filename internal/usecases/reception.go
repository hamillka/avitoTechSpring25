//go:generate mockgen -source=reception.go -destination=./mocks/mock_reception.go -package=mocks
package usecases

import (
	"time"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
)

type ReceptionRepository interface {
	GetLastReception(pvzId string) (models.Reception, error)
	CreateReception(pvzId string) (models.Reception, error)
	UpdateReceptionStatus(recId, status string) (models.Reception, error)
	GetReceptionsByPVZIds(pvzIds []string, startDate, endDate *time.Time) ([]models.Reception, error)
}

type ReceptionService struct {
	pvzRepo PVZRepository
	recRepo ReceptionRepository
}

func NewReceptionService(pvzRepo PVZRepository, recRepo ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		pvzRepo: pvzRepo,
		recRepo: recRepo,
	}
}

func (rs *ReceptionService) CreateReception(pvzId string) (models.Reception, error) {
	/*
		Проверить существует ли такой ПВЗ
		Проверить есть ли активная приемка у этого ПВЗ
		Если есть, то вернуть ошибку
		Если нет, то создать приемку
	*/

	_, err := rs.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	lastReception, err := rs.recRepo.GetLastReception(pvzId)
	if lastReception.Id != "" && lastReception.Status == "in_progress" {
		return models.Reception{}, dto.ErrPVZAlreadyHasReception
	}

	newReception, err := rs.recRepo.CreateReception(pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	return newReception, nil
}
