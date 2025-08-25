package utils

import (
	"time"
	"errors"
	"crypto/rand"
	"encoding/hex"
	"backend-academi/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService maneja la generación y validación de tokens JWT
type JWTService struct {
	secretKey string
	issuer    string
}

// NewJWTService crea una nueva instancia del servicio JWT
func NewJWTService(secretKey, issuer string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

// Claims representa los claims del JWT con información del usuario
type Claims struct {
	UserID    uint            `json:"userId"`
	Email     string          `json:"email"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Role      models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// generateJTI genera un ID único para el JWT
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateToken crea un JWT con los datos del usuario
// Expiración: 15 minutos para tokens de acceso
func (j *JWTService) GenerateToken(user *models.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	
	jti, err := generateJTI()
	if err != nil {
		return "", time.Time{}, err
	}
	
	claims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateToken verifica un JWT y extrae los claims del usuario
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de signing sea HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateRefreshToken genera un token de actualización aleatorio
func (j *JWTService) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsPrivilegedUser verifica si el usuario requiere blacklist (admin o teacher)
func IsPrivilegedUser(role models.UserRole) bool {
	return role == models.AdminRole || role == models.TeacherRole
}