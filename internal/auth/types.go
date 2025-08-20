package auth

import (
	"time"
	"backend-academi/internal/models"
)

// LoginRequest representa los datos de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest representa los datos de registro
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// LoginResponse representa la respuesta después del login
type LoginResponse struct {
	Token     string       `json:"token"`
	User      models.User  `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// TokenResponse para refresh tokens
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}