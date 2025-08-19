package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler maneja las peticiones relacionadas con el estado del servicio
type HealthHandler struct{}

// NewHealthHandler crea una nueva instancia del handler de salud
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthResponse representa la respuesta del endpoint de salud
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version,omitempty"`
}

// CheckHealth verifica el estado del servicio
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Service: "academi-backend",
		Version: "1.0.0", // TODO: Obtener de configuración
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}