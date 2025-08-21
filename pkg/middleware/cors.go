package middleware

import (
	"net/http"
)

// CORSMiddleware maneja la configuración de CORS para conexiones web
type CORSMiddleware struct{}

// NewCORSMiddleware crea una nueva instancia del middleware CORS
func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

// Handler middleware que configura headers CORS
func (m *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Responder a preflight requests (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}