# API Documentation

## Base URL
Por defecto: `http://localhost:8080/api/v1`  
Configurable via variable de entorno `API_BASE_PATH`

**IMPORTANTE para Frontend:** La ruta base debe configurarse en el frontend usando la misma variable de entorno que el backend para mantener consistencia:
```javascript
// Frontend env
REACT_APP_API_BASE_URL=http://localhost:8080/api/v1
// Debe coincidir con API_BASE_PATH del backend
```

## Endpoints Implementados

### Autenticación 

#### POST /auth/register
Registro de nuevo usuario en el sistema.

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "contraseña123",
  "first_name": "Juan",
  "last_name": "Pérez"
}
```

**Validaciones:**
- `email`: Requerido, formato válido de email, único en el sistema
- `password`: Requerido, mínimo 8 caracteres
- `first_name`: Requerido
- `last_name`: Requerido

**Response 201 - Éxito:**
```json
{
  "id": 1,
  "email": "usuario@ejemplo.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "role": "student",
  "is_active": true,
  "created_at": "2025-08-20T14:30:45Z",
  "updated_at": "2025-08-20T14:30:45Z"
}
```

**Response 400 - Error:**
```json
HTTP/1.1 400 Bad Request
Content-Type: text/plain

email already registered
```

---

#### POST /auth/login
Autenticación de usuario existente.

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "contraseña123"
}
```

**Response 200 - Éxito:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "usuario@ejemplo.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "role": "student",
    "is_active": true,
    "created_at": "2025-08-20T14:30:45Z",
    "updated_at": "2025-08-20T14:30:45Z"
  },
  "expires_at": "2025-08-21T14:35:12Z"
}
```

**Response 401 - Error:**
```json
HTTP/1.1 401 Unauthorized
Content-Type: text/plain

invalid credentials
```

---

#### GET /auth/me
Obtener información del usuario autenticado.

**Headers requeridos:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response 200 - Éxito:**
```json
{
  "id": 1,
  "email": "usuario@ejemplo.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "role": "student",
  "is_active": true,
  "created_at": "2025-08-20T14:30:45Z",
  "updated_at": "2025-08-20T14:30:45Z"
}
```

---

### Health Check 

#### GET /health
Verificar estado del servidor.

**Response 200:**
```json
{
  "status": "ok",
  "service": "academi-backend",
  "version": "1.0.0"
}
```

---

## Endpoints Pendientes

### Estudiantes 
[Por completar]

### Cursos 
[Por completar]

### Evaluaciones 
[Por completar]

### Analíticas 
[Por completar]

---

## Autenticación JWT

### Formato del Token
```
Authorization: Bearer <token>
```

### Contenido del JWT (Claims)
```json
{
  "user_id": 1,
  "email": "usuario@ejemplo.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "role": "student",
  "exp": 1724155712,
  "iat": 1724069312,
  "iss": "academi-backend"
}
```

### Expiración
- **Duración**: 24 horas por defecto
- **Configurable**: Via `JWT_EXPIRES_HOURS`

---

## Manejo de Errores

### Códigos de Estado
- **200** - OK (éxito)
- **201** - Created (recurso creado)
- **400** - Bad Request (datos inválidos)
- **401** - Unauthorized (sin autenticación o token inválido)
- **500** - Internal Server Error (error del servidor)

### Formato de Errores
Los errores retornan texto plano:
```
HTTP/1.1 400 Bad Request
Content-Type: text/plain

All fields are required
```