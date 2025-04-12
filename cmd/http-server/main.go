package main

import (
	"net/http"

	"github.com/hamillka/avitoTechSpring25/internal/db"
	"github.com/hamillka/avitoTechSpring25/internal/handlers"
	"github.com/hamillka/avitoTechSpring25/internal/logger"
	"github.com/hamillka/avitoTechSpring25/internal/repositories"
	"github.com/hamillka/avitoTechSpring25/internal/usecases"

	cfg "github.com/hamillka/avitoTechSpring25/internal/config"
)

// @title PVZ Service
// @version 1.0
// @description Avito PVZ Service 2025
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						auth-x
//	@description				Authorization check
func main() {
	config, err := cfg.New()
	logger := logger.CreateLogger(config.Log)

	defer func() {
		err = logger.Sync()
		if err != nil {
			logger.Errorf("Error while syncing logger: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("Something went wrong with config: %v", err)
	}

	db, err := db.CreateConnection(&config.DB)

	defer func() {
		err = db.Close()
		if err != nil {
			logger.Errorf("Error while closing connection to db: %v", err)
		}
	}()

	if err != nil {
		logger.Fatalf("Error while connecting to database: %v", err)
	}

	pr := repositories.NewProductRepository(db)
	pvzr := repositories.NewPVZRepository(db)
	rr := repositories.NewReceptionRepository(db)
	ur := repositories.NewUserRepository(db)

	ps := usecases.NewProductService(pr, rr, pvzr)
	pvzs := usecases.NewPVZService(pvzr, rr, pr)
	rs := usecases.NewReceptionService(pvzr, rr)
	us := usecases.NewUserService(ur)

	r := handlers.Router(ps, pvzs, rs, us, logger)

	port := config.Port
	logger.Info("Server is started on port ", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Fatalf("Error while starting server: %v", err)
	}
}
