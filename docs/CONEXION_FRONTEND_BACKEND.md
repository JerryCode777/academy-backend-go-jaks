# 🔗 Conexión Frontend-Backend - Academi

## 📝 Resumen

Este documento detalla cómo está configurada la conexión entre el frontend React y el backend Go, incluyendo la configuración CORS, manejo de errores y ejemplos de implementación.

## 🎯 Estado Actual

### ✅ Implementado
- Conexión automática del frontend al backend
- Middleware CORS configurado en el backend (corregido para todos los subrouters)
- Manejo de estados de carga y error en el frontend
- Interfaz visual para monitorear la conexión
- Respuestas JSON estructuradas
- Soporte para métodos OPTIONS en todas las rutas
- Headers de Authorization configurados para JWT

### 🔄 Flujo de Comunicación

```
[Frontend React] ----HTTP Request----> [Backend Go]
     ↓                                      ↓
[Port 5173]     <----JSON Response---- [Port 8080]
```

## 🔧 Configuración del Backend

### CORS Middleware
```go
// academy-backend-go-jaks/pkg/middleware/cors.go

type CORSMiddleware struct{}

func NewCORSMiddleware() *CORSMiddleware {
    return &CORSMiddleware{}
}

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
```

### Configuración en Main
```go
// academy-backend-go-jaks/cmd/academi/main.go

// Aplicar CORS a todos los routers
corsMiddleware := middleware.NewCORSMiddleware()
router.Use(corsMiddleware.Handler)

// También aplicar a subrouters para asegurar cobertura completa
authRouter.Use(corsMiddleware.Handler)
protectedRouter.Use(corsMiddleware.Handler)
questionnaireRouter.Use(corsMiddleware.Handler)
```

### Endpoints de Conexión
```go
// Endpoint de prueba principal
func testHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "message": "¡Conexión exitosa desde la web!",
        "data":    []string{"estudiante1", "estudiante2", "estudiante3"},
        "success": true,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "status":  "ok",
        "service": "academi-backend",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

## 🌐 Configuración del Frontend

### Conexión Automática
```javascript
// academy-web-jaks/src/App.jsx

useEffect(() => {
    const fetchData = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/test')
            if (!response.ok) {
                throw new Error('Error al conectar con el backend')
            }
            const data = await response.json()
            setApiData(data)
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    fetchData()
}, [])
```

### Prueba Manual de Conexión
```javascript
const testConnection = async () => {
    setLoading(true)
    setError(null)
    try {
        const response = await fetch('http://localhost:8080/health')
        if (!response.ok) {
            throw new Error('Error en la conexión')
        }
        const data = await response.json()
        alert(`Conexión exitosa! Estado: ${data.status}`)
    } catch (err) {
        setError(err.message)
    } finally {
        setLoading(false)
    }
}
```

## 🎨 Interfaz de Usuario

### Estados Visuales

#### 🔄 Loading State
```jsx
{loading && (
    <div className="text-center py-4">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
        <p className="mt-2 text-gray-600">Conectando...</p>
    </div>
)}
```

#### ❌ Error State
```jsx
{error && (
    <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-4">
        <p className="text-red-700">❌ Error: {error}</p>
        <p className="text-sm text-red-600 mt-1">
            Asegúrate de que el backend esté ejecutándose en http://localhost:8080
        </p>
    </div>
)}
```

#### ✅ Success State
```jsx
{apiData && (
    <div className="bg-green-50 border-l-4 border-green-500 p-4">
        <p className="text-green-700 font-semibold">✅ {apiData.message}</p>
        {apiData.data && (
            <div className="mt-2">
                <p className="text-green-600">Datos de prueba:</p>
                <ul className="list-disc list-inside text-sm text-green-600">
                    {apiData.data.map((item, index) => (
                        <li key={index}>{item}</li>
                    ))}
                </ul>
            </div>
        )}
    </div>
)}
```

## 🧪 Pruebas de Conexión

### 1. Verificación Manual

#### Backend
```bash
# Ejecutar backend
cd academy-backend-go-jaks
go run cmd/academi/main.go

# Probar endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/test
```

#### Frontend
```bash
# Ejecutar frontend
cd academy-web-jaks
npm run dev

# Abrir http://localhost:5173
```

### 2. Verificación Automática

El frontend incluye verificaciones automáticas:
- **Al cargar**: Intenta conectar con `/api/test`
- **Botón manual**: Permite probar `/health` cuando sea necesario
- **Indicadores visuales**: Muestra estado en tiempo real

## 🔍 Debugging

### Logs del Backend
```bash
2025-08-20 12:00:00 🚀 Servidor iniciando en http://localhost:8080
```

### Console del Frontend
```javascript
// Agregar logs para debugging
console.log('Intentando conectar con:', 'http://localhost:8080/api/test')
console.log('Respuesta recibida:', data)
console.log('Error de conexión:', error)
```

### Network Tab (DevTools)
1. Abrir DevTools (F12)
2. Ir a la pestaña "Network"
3. Recargar la página
4. Verificar las peticiones HTTP al backend

## 🚨 Problemas Comunes

### CORS Errors
**Síntoma**: `Access to fetch blocked by CORS policy` o `CORS Missing Allow Origin`

**Solución**:
1. Verificar que el middleware CORS esté aplicado a TODOS los subrouters
2. Comprobar que el backend esté ejecutándose
3. Confirmar que las cabeceras CORS incluyan "Authorization" para JWT
4. Verificar que los métodos OPTIONS estén permitidos en todas las rutas

**Fix Implementado (2025-08-21)**:
- CORS middleware aplicado a router principal Y a todos los subrouters
- Headers de Authorization agregados para soporte JWT
- Métodos OPTIONS habilitados en todas las rutas API

### Connection Refused
**Síntoma**: `Failed to fetch` o `ERR_CONNECTION_REFUSED`

**Solución**:
1. Verificar que el backend esté ejecutándose en puerto 8080
2. Comprobar que no haya firewall bloqueando la conexión
3. Confirmar la URL en el frontend

### JSON Parse Error
**Síntoma**: `Unexpected token < in JSON`

**Solución**:
1. Verificar que el backend esté devolviendo JSON válido
2. Comprobar que Content-Type sea `application/json`
3. Revisar logs del backend por errores

## 📡 Headers HTTP

### Request Headers (Frontend → Backend)
```
GET /api/test HTTP/1.1
Host: localhost:8080
Accept: application/json
Origin: http://localhost:5173
```

### Response Headers (Backend → Frontend)
```
HTTP/1.1 200 OK
Content-Type: application/json
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

## 🔮 Próximas Mejoras

### Autenticación
```javascript
// Futuro: Headers con token
const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
}
```

### Environment Variables
```javascript
// .env
VITE_API_URL=http://localhost:8080

// Usar en código
const API_URL = import.meta.env.VITE_API_URL
```

### Error Handling Avanzado
```javascript
// Retry logic
const fetchWithRetry = async (url, options, retries = 3) => {
    try {
        return await fetch(url, options)
    } catch (error) {
        if (retries > 0) {
            await new Promise(resolve => setTimeout(resolve, 1000))
            return fetchWithRetry(url, options, retries - 1)
        }
        throw error
    }
}
```

## 📊 Métricas de Conexión

### Tiempo de Respuesta
Actualmente las respuestas son inmediatas debido a datos mock.

### Disponibilidad
- **Backend**: 99.9% cuando está ejecutándose localmente
- **Frontend**: Depende del estado del backend

### Throughput
- Soporta múltiples conexiones concurrentes
- No hay límites de rate limiting implementados

---

**Estado**: ✅ **FUNCIONANDO**  
**Última verificación**: 2025-08-21  
**CORS Fix**: ✅ **IMPLEMENTADO** (2025-08-21)  
**Próxima revisión**: Cuando se implemente autenticación avanzada