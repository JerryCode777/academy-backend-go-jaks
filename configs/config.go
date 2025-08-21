package configs

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// DatabaseConfig configuración de base de datos
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string // DEBE venir de env var
	DBName   string
	SSLMode  string
}

// JWTConfig configuración para JWT
type JWTConfig struct {
	SecretKey string // CRÍTICO: debe venir de env var
	Issuer    string
	ExpiresIn int // en horas
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port        int
	Host        string
	APIBasePath string // Ruta base de la API (ej: /api/v1, /academi-api)
}

// LoadConfig carga la configuración desde variables de entorno
// Primero intenta cargar .env, luego usa variables del sistema
func LoadConfig() (*Config, error) {
	// Cargar archivo .env si existe (para desarrollo)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),  // Opcional: default localhost para dev
			Port:     getEnvInt("DB_PORT", 5432),      // Opcional: default Puerto PostgreSQL estándar
			User:     getEnv("DB_USER", "postgres"),   // Opcional: default Usuario PostgreSQL estándar
			Password: getEnv("DB_PASSWORD", ""),       // BLIGATORIO: Sin default por seguridad
			DBName:   getEnv("DB_NAME", "academi"),    // Opcional: default Nombre del proyecto
			SSLMode:  getEnv("DB_SSLMODE", "disable"), // Opcional: disable para dev, require para prod
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", ""),            // OBLIGATORIO: Sin default por seguridad crítica
			Issuer:    getEnv("JWT_ISSUER", "academi-backend"), // Opcional: default Identificador del proyecto
			ExpiresIn: getEnvInt("JWT_EXPIRES_HOURS", 24),      // Opcional: default 24 horas
		},
		Server: ServerConfig{
			Port:        getEnvInt("SERVER_PORT", 8080),     // Opcional: default Puerto estándar desarrollo
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),   // Opcional: default Bind todas las interfaces
			APIBasePath: getEnv("API_BASE_PATH", "/api/v1"), // Opcional: default Estándar versionado
		},
	}

	// Validar configuraciones críticas de seguridad
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate verifica que las configuraciones críticas estén presentes
func (c *Config) Validate() error {
	if c.JWT.SecretKey == "" {
		return errors.New("JWT_SECRET_KEY environment variable is required")
	}

	if len(c.JWT.SecretKey) < 32 {
		return errors.New("JWT_SECRET_KEY must be at least 32 characters long for security")
	}

	if c.Database.Password == "" {
		return errors.New("DB_PASSWORD environment variable is required")
	}

	return nil
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt obtiene una variable de entorno como entero con valor por defecto
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
