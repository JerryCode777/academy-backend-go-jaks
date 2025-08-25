package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware proporciona logging detallado de todas las requests HTTP
// Estilo Django: muestra método, ruta, status, tiempo de respuesta
type LoggingMiddleware struct{}

// NewLoggingMiddleware crea una nueva instancia del middleware de logging
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

// customResponseWriter envuelve el ResponseWriter para capturar el status code y bytes escritos
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (w *customResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *customResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	w.written += len(b)
	return w.ResponseWriter.Write(b)
}

// Handler es el middleware que logea todas las requests HTTP
// Se ejecuta ANTES y DESPUÉS de cada handler
func (l *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// =========================== REQUEST LOGGING ===========================
		log.Printf("[INCOMING] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		// Mostrar query parameters si existen
		if r.URL.RawQuery != "" {
			log.Printf("   Query: %s", r.URL.RawQuery)
		}
		
		// Headers importantes para debugging
		if contentType := r.Header.Get("Content-Type"); contentType != "" {
			log.Printf("   Content-Type: %s", contentType)
		}
		
		if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
			// Acortar User-Agent para mejor legibilidad
			shortUA := userAgent
			if len(shortUA) > 60 {
				shortUA = shortUA[:60] + "..."
			}
			log.Printf("   User-Agent: %s", shortUA)
		}
		
		if auth := r.Header.Get("Authorization"); auth != "" {
			log.Printf("   Authorization: %s", maskAuthToken(auth))
		}

		// Logear el body de requests POST/PUT para debugging
		var bodyBytes []byte
		if r.Method == "POST" || r.Method == "PUT" {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body.Close() // Cerrar el body original
			
			// Crear un nuevo body para que los handlers puedan leerlo
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
			// Mostrar el contenido del body (limitado para seguridad)
			bodyStr := string(bodyBytes)
			if len(bodyStr) > 0 {
				if len(bodyStr) > 500 {
					log.Printf("   BODY: %s... [truncated at 500 chars]", bodyStr[:500])
				} else {
					log.Printf("   BODY: %s", bodyStr)
				}
			} else {
				log.Printf("   BODY: [empty]")
			}
		}

		// Crear wrapper para capturar información de response
		wrapper := &customResponseWriter{
			ResponseWriter: w,
			statusCode:     0,
			written:        0,
		}

		// =========================== EJECUTAR HANDLER ===========================
		next.ServeHTTP(wrapper, r)

		// =========================== RESPONSE LOGGING ===========================
		duration := time.Since(start)
		statusLevel := getStatusLevel(wrapper.statusCode)
		
		log.Printf("[RESPONSE] %s %s %s -> [%d] %s | %d bytes | %v", 
			statusLevel, r.Method, r.URL.Path,
			wrapper.statusCode, http.StatusText(wrapper.statusCode),
			wrapper.written, formatDuration(duration))

		// Alerta para requests lentas
		if duration > 500*time.Millisecond {
			log.Printf("[SLOW] %s %s took %v (>500ms)", r.Method, r.URL.Path, duration)
		}
		
		// Línea separadora para mejor legibilidad
		log.Printf("----------------------------------------------------------------")
	})
}

// getStatusLevel retorna nivel de status para identificación rápida
func getStatusLevel(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "SUCCESS" // 2xx
	case statusCode >= 300 && statusCode < 400:
		return "REDIRECT" // 3xx
	case statusCode >= 400 && statusCode < 500:
		return "CLIENT_ERROR" // 4xx
	case statusCode >= 500:
		return "SERVER_ERROR" // 5xx
	default:
		return "INFO" // 1xx o desconocido
	}
}

// maskAuthToken enmascara el token de autorización para logging seguro
// Muestra solo inicio y final para debugging sin exponer el token completo
func maskAuthToken(auth string) string {
	if len(auth) < 20 {
		return auth // Si es muy corto, mostrarlo completo (probablemente no es JWT)
	}
	
	// Mostrar "Bearer abc...xyz" formato
	return fmt.Sprintf("%s...%s", auth[:12], auth[len(auth)-4:])
}

// formatDuration formatea la duración para mejor legibilidad
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dus", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1000000)
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

// LogServerStartup logea información de inicio del servidor con logging middleware
func LogServerStartup(host string, port int, basePath string) {
	log.Println("")
	log.Println("==================== ACADEMI BACKEND ====================")
	log.Println("[STARTUP] HTTP Server iniciado con logging middleware")
	log.Printf("[CONFIG] Host: %s:%d", host, port)
	log.Printf("[CONFIG] API Base Path: %s", basePath)
	log.Printf("[CONFIG] Full URL: http://%s:%d%s", host, port, basePath)
	log.Println("[INFO] Todas las requests HTTP serán logeadas (estilo Django)")
	log.Println("[FORMAT] [INCOMING] -> [RESPONSE] con métricas detalladas")
	log.Println("===========================================================")
	log.Println("")
}