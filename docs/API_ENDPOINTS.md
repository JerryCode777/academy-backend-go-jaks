# 🔌 API Endpoints - Academi Backend

## 📋 Información General

**Base URL**: `http://localhost:8080`  
**Content-Type**: `application/json`  
**CORS**: Habilitado para todos los orígenes

## 🏠 Endpoints Básicos

### GET /
**Descripción**: Página de inicio de la API

**Respuesta**:
```json
{
  "message": "Bienvenido a Academi Backend API",
  "version": "1.0.0",
  "status": "running"
}
```

**Código de estado**: `200 OK`

---

### GET /health
**Descripción**: Verificar estado del servidor

**Respuesta**:
```json
{
  "status": "ok",
  "service": "academi-backend"
}
```

**Código de estado**: `200 OK`

---

### GET /api/test
**Descripción**: Endpoint de prueba para verificar conexión desde el frontend

**Respuesta**:
```json
{
  "message": "¡Conexión exitosa desde la web!",
  "data": [
    "estudiante1",
    "estudiante2", 
    "estudiante3"
  ],
  "success": true
}
```

**Código de estado**: `200 OK`

## 🔒 CORS Configuration

El servidor está configurado con las siguientes políticas CORS:

- **Access-Control-Allow-Origin**: `*` (todos los orígenes)
- **Access-Control-Allow-Methods**: `GET, POST, OPTIONS`
- **Access-Control-Allow-Headers**: `Content-Type`

### Preflight Requests
El servidor maneja automáticamente las peticiones `OPTIONS` para preflight CORS.

## 📝 Ejemplos de Uso

### cURL Examples

#### Verificar estado del servidor
```bash
curl -X GET http://localhost:8080/health
```

#### Obtener datos de prueba
```bash
curl -X GET http://localhost:8080/api/test
```

#### Página de inicio
```bash
curl -X GET http://localhost:8080/
```

### JavaScript/Frontend Examples

#### Usando fetch()
```javascript
// Verificar conexión
const checkHealth = async () => {
  try {
    const response = await fetch('http://localhost:8080/health');
    const data = await response.json();
    console.log('Server status:', data.status);
  } catch (error) {
    console.error('Connection error:', error);
  }
};

// Obtener datos de prueba
const getTestData = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/test');
    const data = await response.json();
    console.log('Test data:', data.data);
  } catch (error) {
    console.error('Error:', error);
  }
};
```

#### Usando axios
```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080'
});

// Verificar estado
const healthCheck = async () => {
  try {
    const response = await api.get('/health');
    return response.data;
  } catch (error) {
    throw error;
  }
};
```

## 🚀 Cómo Probar los Endpoints

### 1. Asegurar que el servidor esté ejecutándose
```bash
cd academy-backend-go-jaks
go run cmd/academi/main.go
```

### 2. Verificar en el navegador
- Abrir http://localhost:8080 en el navegador
- Debería mostrar el mensaje de bienvenida

### 3. Probar con herramientas HTTP
- **Postman**: Importar las URLs y probar
- **Insomnia**: Crear requests a los endpoints
- **Thunder Client** (VS Code): Extension para probar APIs

### 4. Desde el Frontend React
- Ejecutar `npm run dev` en el directorio `academy-web-jaks`
- La interfaz mostrará automáticamente el estado de la conexión

## 🔧 Configuración del Servidor

### Puerto
El servidor escucha en el puerto **8080** por defecto.

### Middleware
- **CORS**: Configurado para permitir conexiones cross-origin
- **Content-Type**: Todas las respuestas son `application/json`

### Logging
El servidor registra:
- Inicio del servidor: `🚀 Servidor iniciando en http://localhost:8080`
- Errores fatales en caso de fallos

## 📊 Códigos de Estado HTTP

| Código | Descripción |
|--------|-------------|
| 200    | OK - Petición exitosa |
| 404    | Not Found - Endpoint no encontrado |
| 405    | Method Not Allowed - Método HTTP no permitido |
| 500    | Internal Server Error - Error interno del servidor |

## 🔮 Próximos Endpoints (Planificados)

Los siguientes endpoints están en la planificación para futuras versiones:

### Autenticación
- `POST /api/auth/login` - Iniciar sesión
- `POST /api/auth/register` - Registrar usuario
- `POST /api/auth/logout` - Cerrar sesión

### Usuarios
- `GET /api/users/profile` - Obtener perfil de usuario
- `PUT /api/users/profile` - Actualizar perfil

### Cursos
- `GET /api/courses` - Listar cursos
- `GET /api/courses/{id}` - Obtener curso específico
- `POST /api/courses` - Crear curso (admin)

### Evaluaciones
- `GET /api/quizzes` - Listar evaluaciones
- `POST /api/quizzes/{id}/submit` - Enviar respuestas

---

**Última actualización**: 2025-08-20  
**Versión API**: 1.0.0