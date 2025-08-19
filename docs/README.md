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
[Por completar]

## API Endpoints
[Por completar]

## Testing
[Por completar]

## Deployment
[Por completar]
