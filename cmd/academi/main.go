package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"backend-academi/configs"
	"backend-academi/internal/auth"
	"backend-academi/internal/handlers"
	"backend-academi/internal/models"
	"backend-academi/internal/repository"
	"backend-academi/pkg/database"
	"backend-academi/pkg/middleware"
	"backend-academi/pkg/utils"
)

func main() {
	// 1. Cargar configuración desde variables de entorno
	log.Println("Cargando configuración...")
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	log.Printf("Configuración cargada - API Base: %s, Puerto: %d", 
		config.Server.APIBasePath, config.Server.Port)

	// 2. Conectar a la base de datos
	log.Println("Conectando a PostgreSQL...")
	db, err := database.Connect(&config.Database)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	log.Println("Conexión a base de datos establecida")

	// 3. Ejecutar migraciones automáticas
	if err := runMigrations(db); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	// 4. Inicializar servicios (Dependency Injection)
	log.Println("Inicializando servicios...")
	jwtService := utils.NewJWTService(config.JWT.SecretKey, config.JWT.Issuer)
	
	// Repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	blacklistRepo := repository.NewTokenBlacklistRepository(db)
	
	// Services
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, blacklistRepo, jwtService)
	
	// 5. Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	
	// 6. Inicializar middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, jwtService)

	// 7. Configurar router y rutas
	router := mux.NewRouter()
	
	// CORS para permitir conexiones desde web
	router.Use(corsMiddleware)
	
	// Rutas básicas (de conexion-frontend)
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/test", testHandler).Methods("GET")
	
	// Subrouter para API con el base path configurable
	apiRouter := router.PathPrefix(config.Server.APIBasePath).Subrouter()
	
	// Rutas de autenticación (sin middleware)
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	
	// Rutas protegidas (CON middleware de autenticación)
	protectedRouter := apiRouter.PathPrefix("/auth").Subrouter()
	protectedRouter.Use(authMiddleware.RequireAuth)
	protectedRouter.HandleFunc("/me", authHandler.Me).Methods("GET")
	protectedRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// 8. Configurar servidor
	serverAddr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	
	log.Printf("🚀 Servidor iniciando en %s", serverAddr)
	log.Printf("Rutas disponibles:")
	log.Printf("   GET  / - Página de inicio")
	log.Printf("   GET  /health - Health check")
	log.Printf("   GET  /api/test - Endpoint de prueba")
	log.Printf("   POST %s/auth/register", config.Server.APIBasePath)
	log.Printf("   POST %s/auth/login", config.Server.APIBasePath)
	log.Printf("   POST %s/auth/refresh", config.Server.APIBasePath)
	log.Printf("   GET  %s/auth/me (requiere auth)", config.Server.APIBasePath)
	log.Printf("   POST %s/auth/logout (requiere auth)", config.Server.APIBasePath)
	
	log.Fatal(http.ListenAndServe(serverAddr, router))
}

// CORS middleware para conexión web
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

// Página de inicio
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Bienvenido a Academi Backend API",
		"version": "1.0.0",
		"status":  "running",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Verificar estado del servidor  
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"service": "academi-backend",
		"version": "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Endpoint de prueba para la web
func testHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "¡Conexión exitosa desde la web!",
		"data":    []string{"estudiante1", "estudiante2", "estudiante3"},
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// runMigrations ejecuta las migraciones automáticas de GORM
func runMigrations(db *gorm.DB) error {
	log.Println("Ejecutando migraciones de base de datos...")
	
	// GORM creará las tablas automáticamente basándose en los modelos
	return db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.RefreshToken{},
		&models.TokenBlacklist{},
		// TODO: Agregar más modelos cuando se implementen (Course, Quiz, etc.)
	)
}