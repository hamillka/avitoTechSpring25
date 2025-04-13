//go:generate mockgen -source=product.go -destination=./mocks/mock_product.go -package=mocks
package usecases

import (
	"time"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
)

type ProductRepository interface {
	AddProduct(productType, receptionId string) (models.Product, error)
	GetLastProduct(recId string) (models.Product, error)
	DeleteProduct(prodId string) error
	GetProductsByReceptionIds(recIds []string, startDate, endDate *time.Time) ([]models.Product, error)
}

type ProductService struct {
	prodRepo ProductRepository
	recRepo  ReceptionRepository
	pvzRepo  PVZRepository
}

func NewProductService(prodRepo ProductRepository, recRepo ReceptionRepository, pvzRepo PVZRepository) *ProductService {
	return &ProductService{
		prodRepo: prodRepo,
		recRepo:  recRepo,
		pvzRepo:  pvzRepo,
	}
}

func (ps *ProductService) AddProductToReception(productType, pvzId string) (models.Product, error) {
	_, err := ps.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return models.Product{}, err
	}

	lastReception, err := ps.recRepo.GetLastReception(pvzId)
	if err != nil || lastReception.Status != models.INPROGRESS {
		return models.Product{}, dto.ErrNoActiveReception
	}

	product, err := ps.prodRepo.AddProduct(productType, lastReception.Id)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}
