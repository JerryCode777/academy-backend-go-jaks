# Database Configuration

## Configuración de Conexión

### Variables de Entorno
La aplicación usa las siguientes variables para conectarse a PostgreSQL:

```bash
# Obligatorias
DB_PASSWORD=tu_password_postgresql

# Opcionales (con defaults)
DB_HOST=localhost              # default: localhost
DB_PORT=5432                  # default: 5432  
DB_USER=postgres              # default: postgres
DB_NAME=academi               # default: academi
DB_SSLMODE=disable           # default: disable
```

### Configuración Centralizada
La configuración de BD se maneja en `configs/config.go`:

```go
type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string // Sin default por seguridad
    DBName   string
    SSLMode  string
}
```

## Pool de Conexiones

Configuración automática en `pkg/database/connection.go`:

```go
MaxIdleConns: 10        // Conexiones inactivas máximas
MaxOpenConns: 100       // Conexiones totales máximas  
ConnMaxLifetime: 1 hour // Tiempo de vida de conexiones
```

## Migraciones Automáticas

### GORM AutoMigrate
Al arrancar la aplicación, se ejecutan automáticamente:

```go
// En cmd/academi/main.go
db.AutoMigrate(
    &models.User{},
    &models.Student{},
    &models.RefreshToken{},
    &models.TokenBlacklist{},
)
```

**¿Qué hace?**
- Crea tablas si no existen
- Añade nuevas columnas
- Crea índices automáticos
- **NO borra** datos existentes (seguro)

### Logs de Migración
```bash
Ejecutando migraciones de base de datos...
Database connection established successfully
```

## Configuración por Ambiente

### Desarrollo Local
```bash
# .env para desarrollo
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password123
DB_NAME=academi_dev
DB_SSLMODE=disable
```

### Staging
```bash
DB_HOST=staging-db.ejemplo.com
DB_NAME=academi_staging
DB_SSLMODE=require
```

### Producción
```bash
DB_HOST=prod-db.empresa.com
DB_USER=academi_user
DB_PASSWORD=contraseña_super_segura_prod
DB_NAME=academi_prod
DB_SSLMODE=require    # OBLIGATORIO en producción
```

## Configuración con Docker

```bash
# PostgreSQL local con Docker
docker run --name postgres-academi \
  -e POSTGRES_DB=academi_dev \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password123 \
  -p 5432:5432 -d postgres:14
```

## Seguridad Implementada

### Variables de Entorno
- **DB_PASSWORD**: Sin valor por defecto, validación obligatoria
- **Separación**: Diferentes passwords por ambiente
- **SSL**: Forzado en producción via `DB_SSLMODE=require`

### Protección SQL Injection
- **Prepared Statements**: GORM los usa automáticamente
- **Ejemplo seguro**: `db.Where("email = ?", email)`
- **Validación**: En startup si faltan configuraciones críticas

## Esquema de Tablas

### Tablas Actuales

#### refresh_tokens **NUEVA**
Tokens de actualización para el sistema híbrido:
```sql
-- Campos para gestión de sesiones
id, token (unique + index), user_id (FK), expires_at
is_revoked, created_at, updated_at
```

#### token_blacklist **NUEVA**
Blacklist de JWTs para usuarios privilegiados:
```sql
-- Solo para admin/teacher
id, jti (unique + index), token, user_id (FK)
expires_at (index), created_at
```

### Índices Automáticos
- `refresh_tokens.token` - Para búsquedas rápidas de refresh tokens
- `token_blacklist.jti` - Para verificación rápida de blacklist
- `token_blacklist.expires_at` - Para limpieza automática de tokens expirados

### Estrategia de Limpieza
Las tablas nuevas incluyen métodos para limpiar automáticamente registros expirados:
- `RefreshTokenRepository.CleanupExpired()`
- `TokenBlacklistRepository.CleanupExpired()`

## Troubleshooting

### Error de Conexión
```bash
Error connecting to database: failed to connect...
```
**Solución**: Verificar que PostgreSQL esté corriendo y las variables sean correctas

### Error de Migración  
```bash
Error running migrations: ...
```
**Solución**: Verificar permisos del usuario de BD para crear tablas