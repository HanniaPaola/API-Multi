package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"api/paypal/src/archivo/application"
	"api/paypal/src/archivo/core/domain"

	"github.com/gorilla/mux"
)

type PlanHandler struct {
	service *application.PlanService
}

func NewPlanHandler(service *application.PlanService) *PlanHandler {
	return &PlanHandler{service: service}
}

// GetAvailablePlans obtiene los planes disponibles
func (h *PlanHandler) GetAvailablePlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.GetAvailablePlans()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plans)
}

// GetStationPlan obtiene el plan de una estación en particular
func (h *PlanHandler) GetStationPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stationIDStr := vars["stationId"]
	stationID, err := strconv.Atoi(stationIDStr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	station, planDetails, err := h.service.GetStationPlan(stationID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"station": station,
		"plan":    planDetails,
	})
}

// UpgradePlan actualiza el plan de una estación
func (h *PlanHandler) UpgradePlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stationIDStr := vars["stationId"]
	stationID, err := strconv.Atoi(stationIDStr)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	var request struct {
		Plan string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Convertir string a domain.PlanType
	planType := domain.PlanType(request.Plan)
	// Validar que el plan sea válido
	if planType != domain.Basic && planType != domain.BasicPlus {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tipo de plan inválido"})
		return
	}

	updatedStation, err := h.service.UpgradePlan(stationID, planType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStation)
}
