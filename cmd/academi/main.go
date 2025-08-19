package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// main es el punto de entrada de la aplicación Academi
func main() {
	// Inicializar router
	router := mux.NewRouter()
	
	// TODO: Configurar rutas de la API
	router.HandleFunc("/health", healthCheck).Methods("GET")
	
	// TODO: Configurar middleware
	// TODO: Configurar conexión a base de datos
	// TODO: Configurar servicios
	
	log.Println("Servidor iniciando en puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// healthCheck endpoint básico para verificar que el servidor está funcionando
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "academi-backend"}`))
}