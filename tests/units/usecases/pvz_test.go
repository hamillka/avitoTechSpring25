package usecases_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/hamillka/avitoTechSpring25/internal/handlers/dto"
	"github.com/hamillka/avitoTechSpring25/internal/models"
	"github.com/hamillka/avitoTechSpring25/internal/usecases"
	"github.com/hamillka/avitoTechSpring25/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreatePVZ_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mocks.NewMockPVZRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	prodRepo := mocks.NewMockProductRepository(ctrl)

	service := usecases.NewPVZService(pvzRepo, recRepo, prodRepo)

	pvzRepo.EXPECT().CreatePVZ("Москва").Return(models.PVZ{Id: "1", City: "Москва"}, nil)

	pvz, err := service.CreatePVZ("Москва")
	assert.NoError(t, err)
	assert.Equal(t, "Москва", pvz.City)
}

func TestCloseLastReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mocks.NewMockPVZRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	prodRepo := mocks.NewMockProductRepository(ctrl)

	service := usecases.NewPVZService(pvzRepo, recRepo, prodRepo)

	pvzRepo.EXPECT().GetPVZById("pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
	recRepo.EXPECT().GetLastReception("pvz1").Return(models.Reception{Id: "rec1", Status: "in_progress"}, nil)
	recRepo.EXPECT().UpdateReceptionStatus("rec1", "close").Return(models.Reception{Id: "rec1", Status: "close"}, nil)

	rec, err := service.CloseLastReception("pvz1")
	assert.NoError(t, err)
	assert.Equal(t, "close", rec.Status)
}

func TestDeleteLastProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mocks.NewMockPVZRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	prodRepo := mocks.NewMockProductRepository(ctrl)

	service := usecases.NewPVZService(pvzRepo, recRepo, prodRepo)

	pvzRepo.EXPECT().GetPVZById("pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
	recRepo.EXPECT().GetLastReception("pvz1").Return(models.Reception{Id: "rec1", Status: "in_progress"}, nil)
	prodRepo.EXPECT().GetLastProduct("rec1").Return(models.Product{Id: "prod1"}, nil)
	prodRepo.EXPECT().DeleteProduct("prod1").Return(nil)

	err := service.DeleteLastProduct("pvz1")
	assert.NoError(t, err)
}

func TestDeleteLastProduct_NoProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mocks.NewMockPVZRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	prodRepo := mocks.NewMockProductRepository(ctrl)

	service := usecases.NewPVZService(pvzRepo, recRepo, prodRepo)

	pvzRepo.EXPECT().GetPVZById("pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
	recRepo.EXPECT().GetLastReception("pvz1").Return(models.Reception{Id: "rec1", Status: "in_progress"}, nil)
	prodRepo.EXPECT().GetLastProduct("rec1").Return(models.Product{}, dto.ErrNoProductsInReception)

	err := service.DeleteLastProduct("pvz1")
	assert.ErrorIs(t, err, dto.ErrNoProductsInReception)
}

func TestGetPVZWithPagination_EmptyDates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPVZRepo := mocks.NewMockPVZRepository(ctrl)
	mockRecRepo := mocks.NewMockReceptionRepository(ctrl)
	mockProdRepo := mocks.NewMockProductRepository(ctrl)

	pvzs := []models.PVZ{
		{
			Id:               "pvz1",
			RegistrationDate: time.Now().String(),
			City:             "Moscow",
		},
	}

	receptions := []models.Reception{
		{
			Id:       "rec1",
			DateTime: time.Now().String(),
			PVZId:    "pvz1",
			Status:   "open",
		},
	}

	products := []models.Product{
		{
			Id:          "prod1",
			DateTime:    time.Now().String(),
			Type:        "phone",
			ReceptionId: "rec1",
		},
	}

	mockPVZRepo.EXPECT().GetPVZsWithPagination(0, 10).Return(pvzs, nil)

	mockRecRepo.EXPECT().GetReceptionsByPVZIds([]string{"pvz1"}, nil, nil).Return(receptions, nil)

	mockProdRepo.EXPECT().GetProductsByReceptionIds([]string{"rec1"}, nil, nil).Return(products, nil)

	service := usecases.NewPVZService(mockPVZRepo, mockRecRepo, mockProdRepo)

	result, err := service.GetPVZWithPagination(nil, nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "pvz1", result[0].PVZ.Id)
	assert.Equal(t, 1, len(result[0].Receptions))
	assert.Equal(t, "rec1", result[0].Receptions[0].Reception.Id)
	assert.Equal(t, 1, len(result[0].Receptions[0].Products))
	assert.Equal(t, "prod1", result[0].Receptions[0].Products[0].Id)
}
