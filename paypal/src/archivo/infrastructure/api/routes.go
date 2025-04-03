package api

import (
	"api/paypal/src/archivo/application"
	"api/paypal/src/archivo/infrastructure/api/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, planService *application.PlanService) {
	planHandler := handlers.NewPlanHandler(planService)

	r.HandleFunc("/plans", planHandler.GetAvailablePlans).Methods("GET")
	r.HandleFunc("/stations/{stationId}/plan", planHandler.GetStationPlan).Methods("GET")
	r.HandleFunc("/stations/{stationId}/plan", planHandler.UpgradePlan).Methods("PATCH")
}
