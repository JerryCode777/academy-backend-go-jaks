package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	
	// CORS para permitir conexiones desde web
	router.Use(corsMiddleware)
	
	// Rutas básicas
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/test", testHandler).Methods("GET")
	
	log.Println("🚀 Servidor iniciando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// CORS middleware para conexión web
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Página de inicio
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Bienvenido a Academi Backend API",
		"version": "1.0.0",
		"status":  "running",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Verificar estado del servidor
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"service": "academi-backend",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Endpoint de prueba para la web
func testHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "¡Conexión exitosa desde la web!",
		"data":    []string{"estudiante1", "estudiante2", "estudiante3"},
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}