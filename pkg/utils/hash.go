package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashService maneja el hashing y verificación de contraseñas
type HashService struct{}

// NewHashService crea una nueva instancia del servicio de hash
func NewHashService() *HashService {
	return &HashService{}
}

// HashPassword convierte una contraseña en texto plano a un hash seguro usando bcrypt
// Cost 12 proporciona un balance entre seguridad y performance en 2025
func (h *HashService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword compara una contraseña en texto plano con su hash
// Retorna true si coinciden, false si no
func (h *HashService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}