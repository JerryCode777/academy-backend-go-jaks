# Academi Backend

## Descripción
Backend para Academi, una aplicación inteligente diseñada para ayudar a estudiantes preuniversitarios a lograr sus objetivos académicos.

## Equipo de Desarrollo
Este proyecto será desarrollado por un equipo de 2 desarrolladoras utilizando Go como lenguaje principal.

## Arquitectura
- **Lenguaje**: Go
- **Base de datos**: PostgreSQL (recomendado)
- **Arquitectura**: Clean Architecture / Hexagonal Architecture
- **API**: RESTful con posibilidad de GraphQL

## Estructura del Proyecto - Esto puede ser modificado
```
backend-academi/
├── cmd/academi/           # Punto de entrada de la aplicación
├── internal/              # Código interno de la aplicación
│   ├── auth/             # Autenticación y autorización
│   ├── student/          # Gestión de estudiantes
│   ├── course/           # Gestión de cursos
│   ├── quiz/             # Sistema de evaluaciones
│   ├── analytics/        # Análisis y métricas
│   ├── handlers/         # Controladores HTTP
│   ├── models/           # Modelos de datos
│   ├── services/         # Lógica de negocio
│   └── repository/       # Acceso a datos
├── pkg/                  # Paquetes reutilizables
│   ├── database/         # Configuración de base de datos
│   ├── utils/            # Utilidades generales
│   └── middleware/       # Middleware HTTP
├── configs/              # Archivos de configuración
└── docs/                 # Documentación del proyecto
```

## Instalación y Configuración

### Prerequisitos
- Go 1.21+
- PostgreSQL 14+

### Pasos básicos
```bash
# 1. Instalar dependencias
make deps

# 2. Configurar variables de entorno (ver DEPLOYMENT.md)
cp .env.dev .env
# Editar .env con tus configuraciones

# 3. Ejecutar aplicación
make run
```

## Comandos Principales

```bash
make deps          # Descargar dependencias
make run           # Ejecutar aplicación
make build         # Compilar binario
make test          # Ejecutar tests
make clean         # Limpiar archivos generados
```

## Documentación

- **[API.md](./API.md)** - Endpoints, requests/responses, autenticación
- **[DATABASE.md](./DATABASE.md)** - Configuración BD, modelos, migraciones  
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - Variables de entorno, configuración producción

## Estado Actual

### Funcionalidades Implementadas
- Sistema de autenticación completo (registro, login, JWT, refresh tokens)
- Logout híbrido (elimina refresh tokens + blacklist para usuarios privilegiados)
- Middleware de autenticación para rutas protegidas
- **Conexión con frontend habilitada** (CORS configurado)
- Endpoints básicos de prueba para validar conexión
- Configuración centralizada con variables de entorno
- Migraciones automáticas de base de datos
- Arquitectura modularizada y escalable

### Endpoints de Conexión Frontend
- `GET /` - Página de bienvenida
- `GET /health` - Health check para monitoreo
- `GET {API_BASE_PATH}/test` - Endpoint de prueba para validar conexión frontend
- Todos los endpoints auth funcionan con CORS habilitado

### CORS y Frontend
✅ **Conexión frontend lista**: El backend está configurado con CORS para permitir conexiones desde aplicaciones web.

**Headers CORS configurados:**
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, OPTIONS`  
- `Access-Control-Allow-Headers: Content-Type, Authorization`
