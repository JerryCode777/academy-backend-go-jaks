package main

import (
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
	"backend-academi/internal/services"
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
	questionnaireRepo := repository.NewQuestionnaireRepository(db)
	
	// Services
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, blacklistRepo, jwtService)
	questionnaireService := services.NewQuestionnaireService(questionnaireRepo, userRepo)
	
	// 5. Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	basicHandler := handlers.NewBasicHandler()
	questionnaireHandler := handlers.NewQuestionnaireHandler(questionnaireService)
	
	// 6. Inicializar middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, jwtService)
	corsMiddleware := middleware.NewCORSMiddleware()
	loggingMiddleware := middleware.NewLoggingMiddleware()

	// 7. Configurar router y rutas
	router := mux.NewRouter()
	
	// CORS para permitir conexiones desde web
	router.Use(corsMiddleware.Handler)
	
	// Logging middleware para ver todas las requests HTTP (estilo Django)
	router.Use(loggingMiddleware.Handler)
	
	// Rutas básicas (de conexion-frontend)
	router.HandleFunc("/", basicHandler.Home).Methods("GET")
	router.HandleFunc("/health", basicHandler.Health).Methods("GET")
	
	// Subrouter para API con el base path configurable
	apiRouter := router.PathPrefix(config.Server.APIBasePath).Subrouter()
	
	// Rutas API básicas (usando API_BASE_PATH)
	apiRouter.HandleFunc("/test", basicHandler.Test).Methods("GET")
	
	// Rutas de autenticación (sin middleware)
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.Use(corsMiddleware.Handler) // Aplicar CORS también a subrouters
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST", "OPTIONS")
	
	// Rutas protegidas (CON middleware de autenticación)
	protectedRouter := apiRouter.PathPrefix("/auth").Subrouter()
	protectedRouter.Use(corsMiddleware.Handler) // Aplicar CORS
	protectedRouter.Use(authMiddleware.RequireAuth)
	protectedRouter.HandleFunc("/me", authHandler.Me).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")
	
	// Rutas de cuestionarios públicas (sin autenticación)
	questionnairePublicRouter := apiRouter.PathPrefix("/questionnaire").Subrouter()
	questionnairePublicRouter.Use(corsMiddleware.Handler) // Aplicar CORS
	questionnairePublicRouter.HandleFunc("/initial/public", questionnaireHandler.GetInitialQuestionnairePublic).Methods("GET", "OPTIONS")
	
	// Rutas de cuestionarios protegidas (CON middleware de autenticación)
	questionnaireRouter := apiRouter.PathPrefix("/questionnaire").Subrouter()
	questionnaireRouter.Use(corsMiddleware.Handler) // Aplicar CORS
	questionnaireRouter.Use(authMiddleware.RequireAuth)
	questionnaireRouter.HandleFunc("/initial", questionnaireHandler.GetInitialQuestionnaire).Methods("GET", "OPTIONS")
	questionnaireRouter.HandleFunc("/initial/submit", questionnaireHandler.SubmitInitialQuestionnaire).Methods("POST", "OPTIONS")
	questionnaireRouter.HandleFunc("/initial/status", questionnaireHandler.CheckInitialCompletion).Methods("GET", "OPTIONS")
	questionnaireRouter.HandleFunc("/initial/response", questionnaireHandler.GetUserInitialResponse).Methods("GET", "OPTIONS")

	// 8. Configurar servidor
	serverAddr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	
	// Usar la función de logging para startup más limpio
	middleware.LogServerStartup(config.Server.Host, config.Server.Port, config.Server.APIBasePath)
	
	log.Printf("[ROUTES] Endpoints disponibles:")
	log.Printf("   GET  / - Página de inicio")
	log.Printf("   GET  /health - Health check")
	log.Printf("   GET  %s/test - Endpoint de prueba", config.Server.APIBasePath)
	log.Printf("   POST %s/auth/register", config.Server.APIBasePath)
	log.Printf("   POST %s/auth/login", config.Server.APIBasePath)
	log.Printf("   POST %s/auth/refresh", config.Server.APIBasePath)
	log.Printf("   GET  %s/auth/me (requiere auth)", config.Server.APIBasePath)
	log.Printf("   POST %s/auth/logout (requiere auth)", config.Server.APIBasePath)
	log.Printf("   GET  %s/questionnaire/initial/public - Cuestionario inicial público", config.Server.APIBasePath)
	log.Printf("   GET  %s/questionnaire/initial (requiere auth)", config.Server.APIBasePath)
	log.Printf("   POST %s/questionnaire/initial/submit (requiere auth)", config.Server.APIBasePath)
	log.Printf("   GET  %s/questionnaire/initial/status (requiere auth)", config.Server.APIBasePath)
	log.Printf("   GET  %s/questionnaire/initial/response (requiere auth)", config.Server.APIBasePath)
	
	log.Fatal(http.ListenAndServe(serverAddr, router))
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
		&models.Questionnaire{},
		&models.QuestionnaireResponse{},
		// TODO: Agregar más modelos cuando se implementen (Course, Quiz, etc.)
	)
}