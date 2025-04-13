package usecases_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/hamillka/avitoTechSpring25/internal/usecases"
	"github.com/hamillka/avitoTechSpring25/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddProductToReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mocks.NewMockProductRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	pvzRepo := mocks.NewMockPVZRepository(ctrl)

	service := usecases.NewProductService(prodRepo, recRepo, pvzRepo)

	pvzRepo.EXPECT().GetPVZById("pvz123").Return(models.PVZ{Id: "pvz123"}, nil)
	recRepo.EXPECT().GetLastReception("pvz123").Return(models.Reception{Id: "rec1", Status: "in_progress"}, nil)
	prodRepo.EXPECT().AddProduct("type1", "rec1").Return(models.Product{Id: "prod1", Type: "type1"}, nil)

	product, err := service.AddProductToReception("type1", "pvz123")

	require.NoError(t, err)
	assert.Equal(t, "prod1", product.Id)
}

func TestAddProductToReception_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mocks.NewMockProductRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	pvzRepo := mocks.NewMockPVZRepository(ctrl)

	service := usecases.NewProductService(prodRepo, recRepo, pvzRepo)

	pvzRepo.EXPECT().GetPVZById("pvz404").Return(models.PVZ{}, errors.New("not found"))

	_, err := service.AddProductToReception("type1", "pvz404")

	assert.Error(t, err)
}

func TestAddProductToReception_NoActiveReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mocks.NewMockProductRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	pvzRepo := mocks.NewMockPVZRepository(ctrl)

	service := usecases.NewProductService(prodRepo, recRepo, pvzRepo)

	pvzRepo.EXPECT().GetPVZById("pvz123").Return(models.PVZ{Id: "pvz123"}, nil)
	recRepo.EXPECT().GetLastReception("pvz123").Return(models.Reception{Id: "rec1", Status: "close"}, nil)

	_, err := service.AddProductToReception("type1", "pvz123")

	assert.ErrorIs(t, err, dto.ErrNoActiveReception)
}
