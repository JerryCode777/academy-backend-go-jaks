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
  "refresh_token": "a1b2c3d4e5f6...",
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
  "expires_at": "2025-08-21T14:50:12Z"
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

#### POST /auth/refresh
Renovar token de acceso usando refresh token.

**Request Body:**
```json
{
  "refresh_token": "a1b2c3d4e5f6..."
}
```

**Response 200 - Éxito:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6...",
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
  "expires_at": "2025-08-21T15:05:30Z"
}
```

**Response 401 - Error:**
```json
HTTP/1.1 401 Unauthorized
Content-Type: text/plain

invalid or expired refresh token
```

---

#### POST /auth/logout
Cerrar sesión del usuario (requiere autenticación).

**Headers requeridos:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response 200 - Éxito:**
```json
{
  "message": "Successfully logged out"
}
```

**Comportamiento por tipo de usuario:**
- **Estudiantes**: Elimina refresh token, JWT expira naturalmente en 15 minutos
- **Admin/Teacher**: Elimina refresh token Y añade JWT a blacklist para invalidación inmediata

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

## Sistema de Autenticación Híbrido

### Formato del Token
```
Authorization: Bearer <token>
```

### Contenido del JWT (Claims)
```json
{
  "jti": "abc123def456...",
  "user_id": 1,
  "email": "usuario@ejemplo.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "role": "student",
  "exp": 1724070212,
  "iat": 1724069312,
  "iss": "academi-backend"
}
```

### Duración de Tokens
- **Access Token (JWT)**: 15 minutos
- **Refresh Token**: 7 días

### Estrategia por Tipo de Usuario

#### 🔵 Usuarios Normales (student)
- **Access Token**: 15 minutos
- **Refresh Token**: Almacenado en BD, 7 días de duración
- **Logout**: Solo elimina refresh token de BD
- **Seguridad**: JWT expira naturalmente en máximo 15 minutos

#### 🔴 Usuarios Privilegiados (admin, teacher)
- **Access Token**: 15 minutos + JTI único para tracking
- **Refresh Token**: Almacenado en BD, 7 días de duración
- **Blacklist**: JWT añadido a blacklist en BD al hacer logout
- **Logout**: Elimina refresh token + añade JWT a blacklist
- **Seguridad**: Invalidación inmediata + verificación en cada request

### Flujo de Renovación
1. Access token expira (15 min)
2. Frontend usa refresh token para obtener nuevo access token
3. Refresh token se puede reutilizar hasta su expiración (7 días)
4. Si refresh token expira, usuario debe hacer login nuevamente

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