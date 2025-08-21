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
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	User         models.User  `json:"user"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// RefreshTokenRequest representa la solicitud de refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}