package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"api/consumer/src/application"
)

type SubscriptionHandler struct {
	notificationAppService *application.NotificationAppService
}

func NewSubscriptionHandler(notificationAppService *application.NotificationAppService) *SubscriptionHandler {
	return &SubscriptionHandler{notificationAppService: notificationAppService}
}

func (h *SubscriptionHandler) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error": "Solicitud inválida"}`, http.StatusBadRequest)
		fmt.Println("Error decodificando solicitud:", err)
		return
	}

	if req.Token == "" {
		http.Error(w, `{"error": "Token requerido"}`, http.StatusBadRequest)
		fmt.Println("Error: Token vacío")
		return
	}

	id, err := h.notificationAppService.Subscribe(req.Token)
	if err != nil {
		http.Error(w, `{"error": "Error guardando en BD"}`, http.StatusInternalServerError)
		fmt.Println("Error guardando en BD:", err)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"userId":  id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}