# Deployment Guide

## Variables de Entorno

### Obligatorias (SIN defaults)
Estas variables DEBEN configurarse o la aplicación no arrancará:

```bash
# Base de datos - Contraseña de PostgreSQL
DB_PASSWORD=tu_password_postgresql_seguro

# JWT - Clave secreta para firmar tokens (mínimo 32 caracteres)
JWT_SECRET_KEY=tu_clave_super_secreta_de_minimo_32_caracteres
```

### Opcionales (CON defaults)
Si no se configuran, usan valores por defecto:

```bash
# === BASE DE DATOS ===
DB_HOST=localhost              # default: localhost
DB_PORT=5432                  # default: 5432
DB_USER=postgres              # default: postgres
DB_NAME=academi               # default: academi
DB_SSLMODE=disable           # default: disable (dev), require (prod)

# === JWT CONFIGURACIÓN ===
JWT_ISSUER=academi-backend    # default: academi-backend
JWT_EXPIRES_HOURS=24         # default: 24

# === SERVIDOR ===
SERVER_PORT=8080             # default: 8080
SERVER_HOST=0.0.0.0          # default: 0.0.0.0
API_BASE_PATH=/api/v1        # default: /api/v1
```

## Desarrollo Local

### 1. Configurar archivo .env
```bash
# Copiar template de desarrollo
cp .env.dev .env

# Editar con tus valores
nano .env
```

### 2. Configurar PostgreSQL local
```bash
# Crear base de datos
createdb academi_dev

# O usar Docker
docker run --name postgres-academi \
  -e POSTGRES_DB=academi_dev \
  -e POSTGRES_PASSWORD=password123 \
  -p 5432:5432 -d postgres:14
```

### 3. Ejecutar aplicación
```bash
make run
# o
go run ./cmd/academi
```

## Staging

### Variables específicas de staging:
```bash
# Base de datos staging
DB_HOST=staging-db.ejemplo.com
DB_NAME=academi_staging
DB_SSLMODE=require

# JWT con clave diferente a producción
JWT_SECRET_KEY=clave_diferente_para_staging

# Puerto específico
SERVER_PORT=8080
API_BASE_PATH=/api/v1
```

## Producción

### Configuración recomendada:
```bash
# Base de datos producción
DB_HOST=prod-db.ejemplo.com
DB_NAME=academi_prod
DB_SSLMODE=require
DB_PASSWORD=contraseña_super_segura_produccion

# JWT producción
JWT_SECRET_KEY=clave_extremadamente_segura_produccion_64_caracteres_minimo
JWT_EXPIRES_HOURS=24

# Servidor producción
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
API_BASE_PATH=/api/v1
```

### Consideraciones de seguridad:
- **JWT_SECRET_KEY**: Genera una clave criptográficamente segura
- **DB_PASSWORD**: Usa contraseñas fuertes y diferentes por ambiente
- **DB_SSLMODE**: SIEMPRE usar "require" en producción
- **Firewall**: Solo puerto necesario abierto
- **HTTPS**: Usar proxy reverso (nginx) con SSL

## Frontend Configuración

El frontend debe configurar la URL base usando el mismo path que el backend:

```javascript
// React (.env)
REACT_APP_API_BASE_URL=http://localhost:8080/api/v1

// Vue (.env)
VUE_APP_API_BASE_URL=http://localhost:8080/api/v1

// Next.js (.env.local)
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

**IMPORTANTE:** Si cambias `API_BASE_PATH` en el backend, también cámbialo en el frontend.

## Validación de Configuración

La aplicación validará automáticamente al arrancar:
- JWT_SECRET_KEY debe tener mínimo 32 caracteres
- DB_PASSWORD no puede estar vacío
- Si falta alguna variable obligatoria, la app NO arrancará

```bash
# Ejemplo de error de validación:
Error loading configuration: JWT_SECRET_KEY must be at least 32 characters long for security
```

## Docker (Pendiente)
[Por implementar en futuras versiones]