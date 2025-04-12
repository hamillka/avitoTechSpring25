package main

import (
	"fmt"
	"net"

	"github.com/hamillka/avitoTechSpring25/internal/config"
	"github.com/hamillka/avitoTechSpring25/internal/db"

	mygrpc "github.com/hamillka/avitoTechSpring25/internal/grpc"
	"github.com/hamillka/avitoTechSpring25/internal/grpc/pvz_v1"
	"github.com/hamillka/avitoTechSpring25/internal/logger"
	"github.com/hamillka/avitoTechSpring25/internal/repositories"
	"github.com/hamillka/avitoTechSpring25/internal/usecases"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.New()
	logger := logger.CreateLogger(cfg.Log)

	defer func() {
		err = logger.Sync()
		if err != nil {
			logger.Errorf("Error while syncing logger: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("Something went wrong with config: %v", err)
	}

	db, err := db.CreateConnection(&cfg.DB)

	defer func() {
		err = db.Close()
		if err != nil {
			logger.Errorf("Error while closing connection to db: %v", err)
		}
	}()

	if err != nil {
		logger.Fatalf("Error while connecting to database: %v", err)
	}

	pvzRepo := repositories.NewPVZRepository(db)
	recRepo := repositories.NewReceptionRepository(db)
	prodRepo := repositories.NewProductRepository(db)
	pvzService := usecases.NewPVZService(pvzRepo, recRepo, prodRepo)

	srv := grpc.NewServer()
	pvz_v1.RegisterPVZServiceServer(srv, mygrpc.NewPVZServer(pvzService))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	logger.Info("Starting gRPC server on port ", cfg.GRPCPort)
	if err := srv.Serve(lis); err != nil {
		logger.Fatalf("gRPC server failed: %v", err)
	}
}
