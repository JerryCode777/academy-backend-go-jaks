package handlers

import (
	"encoding/json"
	"net/http"
)

// BasicHandler maneja los endpoints básicos de la aplicación
type BasicHandler struct{}

// NewBasicHandler crea una nueva instancia del handler básico
func NewBasicHandler() *BasicHandler {
	return &BasicHandler{}
}

// Home maneja la página de inicio GET /
func (h *BasicHandler) Home(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Bienvenido a Academi Backend API",
		"version": "1.0.0",
		"status":  "running",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health maneja el endpoint de health check GET /health
func (h *BasicHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"service": "academi-backend",
		"version": "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Test maneja el endpoint de prueba para frontend GET {API_BASE_PATH}/test
func (h *BasicHandler) Test(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "¡Conexión exitosa desde la web!",
		"data":    []string{"estudiante1", "estudiante2", "estudiante3"},
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}