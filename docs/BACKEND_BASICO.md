# Backend Básico - Academi API

## ¿Qué hemos creado?

Un servidor web básico en Go que puede conectarse con aplicaciones web (frontend).

## Librerías utilizadas

### 1. Gorilla Mux (`github.com/gorilla/mux`)
**¿Por qué?** Es el router más popular de Go para APIs REST.
**¿Qué hace?** Maneja las rutas URL (ej: `/health`, `/api/test`)

```go
router := mux.NewRouter()
router.HandleFunc("/health", healthHandler).Methods("GET")
```

### 2. Encoding/JSON (librería estándar de Go)
**¿Por qué?** Para convertir datos de Go a JSON que entiende la web.
**¿Qué hace?** Transforma `map[string]string` a `{"key": "value"}`

```go
json.NewEncoder(w).Encode(response)
```

### 3. Net/HTTP (librería estándar de Go)
**¿Por qué?** Es la base de todos los servidores web en Go.
**¿Qué hace?** Crea el servidor que escucha en el puerto 8080.

## Endpoints creados

### 1. `/` - Página de inicio
- **Método:** GET
- **Respuesta:** Información básica de la API
- **Uso:** Verificar que el servidor está funcionando

### 2. `/health` - Estado del servidor
- **Método:** GET  
- **Respuesta:** `{"status": "ok", "service": "academi-backend"}`
- **Uso:** Monitoreo del servidor

### 3. `/api/test` - Prueba de conexión web
- **Método:** GET
- **Respuesta:** Datos de prueba para verificar conexión
- **Uso:** Probar desde aplicaciones web

## CORS - Conexión Web

### ¿Qué es CORS?
CORS (Cross-Origin Resource Sharing) permite que una web en `http://localhost:3000` se conecte a nuestro backend en `http://localhost:8080`.

### ¿Por qué lo necesitamos?
Sin CORS, los navegadores bloquean las conexiones entre diferentes puertos/dominios por seguridad.

### Código CORS:
```go
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
```

## Cómo probar

### 1. Ejecutar el servidor:
```bash
go run cmd/academi/main.go
```

### 2. Probar endpoints:
```bash
# Página de inicio
curl http://localhost:8080/

# Estado del servidor
curl http://localhost:8080/health

# Prueba de API
curl http://localhost:8080/api/test
```

### 3. Desde JavaScript (web):
```javascript
// Ejemplo para conectar desde una web
fetch('http://localhost:8080/api/test')
  .then(response => response.json())
  .then(data => console.log(data));
```

## Estructura del código

```
cmd/academi/main.go
├── main() - Punto de entrada
├── corsMiddleware() - Permite conexiones web
├── homeHandler() - Maneja "/"
├── healthHandler() - Maneja "/health"  
└── testHandler() - Maneja "/api/test"
```

## Próximos pasos

1. ✅ Servidor básico funcionando
2. ✅ CORS configurado para web
3. ✅ Endpoints de prueba
4. 🔄 Conectar con base de datos 
5. 🔄 Añadir más endpoints según necesidades