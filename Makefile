# Variables
BINARY_NAME=academi
MAIN_PATH=./cmd/academi
BUILD_DIR=./build

# Comandos por defecto
.PHONY: help build run clean test deps dev

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Construir el binario
	@echo "Construyendo $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: ## Ejecutar la aplicación
	@echo "Ejecutando $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

clean: ## Limpiar archivos generados
	@echo "Limpiando archivos..."
	@rm -rf $(BUILD_DIR)
	@go clean

test: ## Ejecutar tests
	@echo "Ejecutando tests..."
	@go test -v ./...

test-coverage: ## Ejecutar tests con coverage
	@echo "Ejecutando tests con coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

deps: ## Descargar dependencias
	@echo "Descargando dependencias..."
	@go mod download
	@go mod tidy

dev: ## Modo desarrollo con hot reload (requiere air)
	@echo "Iniciando en modo desarrollo..."
	@air

fmt: ## Formatear código
	@echo "Formateando código..."
	@go fmt ./...

lint: ## Ejecutar linter
	@echo "Ejecutando linter..."
	@golangci-lint run

install-tools: ## Instalar herramientas de desarrollo
	@echo "Instalando herramientas..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest