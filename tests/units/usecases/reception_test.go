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
)

func TestCreateReception_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mocks.NewMockPVZRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)

	service := usecases.NewReceptionService(pvzRepo, recRepo)

	pvzRepo.EXPECT().GetPVZById("pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
	recRepo.EXPECT().GetLastReception("pvz1").Return(models.Reception{Status: "close"}, nil)
	recRepo.EXPECT().CreateReception("pvz1").Return(models.Reception{Id: "rec1", Status: "in_progress"}, nil)

	reception, err := service.CreateReception("pvz1")

	assert.NoError(t, err)
	assert.Equal(t, "rec1", reception.Id)
	assert.Equal(t, "in_progress", reception.Status)
}

func TestCreateReception_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mocks.NewMockPVZRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)

	service := usecases.NewReceptionService(pvzRepo, recRepo)

	pvzRepo.EXPECT().GetPVZById("unknown").Return(models.PVZ{}, errors.New("not found"))

	_, err := service.CreateReception("unknown")

	assert.Error(t, err)
}

func TestCreateReception_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mocks.NewMockPVZRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)

	service := usecases.NewReceptionService(pvzRepo, recRepo)

	pvzRepo.EXPECT().GetPVZById("pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
	recRepo.EXPECT().GetLastReception("pvz1").Return(models.Reception{Id: "rec1", Status: "in_progress"}, nil)

	_, err := service.CreateReception("pvz1")

	assert.ErrorIs(t, err, dto.ErrPVZAlreadyHasReception)
}
