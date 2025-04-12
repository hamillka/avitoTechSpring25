package handlers

import (
	"github.com/gorilla/mux"
	_ "github.com/hamillka/avitoTechSpring25/api"
	"github.com/hamillka/avitoTechSpring25/internal/handlers/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func Router(
	ps ProductService,
	pvzs PVZService,
	rs ReceptionService,
	us UserService,
	logger *zap.SugaredLogger,
) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.MetricsMiddleware)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	auth := router.PathPrefix("").Subrouter()
	fun := router.PathPrefix("").Subrouter()

	fun.Use(middlewares.AuthMiddleware)

	ph := NewProductHandler(ps, logger)
	pvzh := NewPVZHandler(pvzs, logger)
	rh := NewReceptionHandler(rs, logger)
	uh := NewUserHandler(us, logger)

	auth.HandleFunc("/login", uh.Login).Methods("POST")
	auth.HandleFunc("/register", uh.Register).Methods("POST")
	auth.HandleFunc("/dummyLogin", uh.DummyLogin).Methods("POST")

	fun.HandleFunc("/pvz", pvzh.CreatePVZ).Methods("POST")
	fun.HandleFunc("/pvz", pvzh.GetPVZWithPagination).Methods("GET")
	fun.HandleFunc("/pvz/{pvzId}/close_last_reception", pvzh.CloseLastReception).Methods("POST")
	fun.HandleFunc("/pvz/{pvzId}/delete_last_product", pvzh.DeleteLastProduct).Methods("POST")

	fun.HandleFunc("/receptions", rh.CreateReception).Methods("POST")
	fun.HandleFunc("/products", ph.AddProductToReception).Methods("POST")

	return router
}
