package grpc

import (
	"context"
	"time"

	pvz_v1 "github.com/hamillka/avitoTechSpring25/internal/grpc/pvz_v1"
	"github.com/hamillka/avitoTechSpring25/internal/usecases"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZServer struct {
	pvz_v1.UnimplementedPVZServiceServer
	service *usecases.PVZService
}

func NewPVZServer(s *usecases.PVZService) *PVZServer {
	return &PVZServer{service: s}
}

func (s *PVZServer) GetPVZList(ctx context.Context, req *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	pvzs, err := s.service.GetAllPVZs(ctx)
	if err != nil {
		return nil, err
	}

	response := &pvz_v1.GetPVZListResponse{}
	for _, p := range pvzs {
		regDate, _ := time.Parse(time.RFC3339, p.RegistrationDate)
		response.Pvzs = append(response.Pvzs, &pvz_v1.PVZ{
			Id:               p.Id,
			RegistrationDate: timestamppb.New(regDate),
			City:             p.City,
		})
	}
	return response, nil
}
