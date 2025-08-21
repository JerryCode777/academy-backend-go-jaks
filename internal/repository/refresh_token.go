package repository

import (
	"backend-academi/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// RefreshTokenRepository maneja todas las operaciones de base de datos para refresh tokens
type RefreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository crea una nueva instancia del repositorio de refresh tokens
func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create inserta un nuevo refresh token en la base de datos
func (r *RefreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	result := r.db.Create(refreshToken)
	return result.Error
}

// GetByToken busca un refresh token por su valor
func (r *RefreshTokenRepository) GetByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Preload("User").Where("token = ? AND is_revoked = false AND expires_at > ?", token, time.Now()).First(&refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found or expired")
		}
		return nil, err
	}
	return &refreshToken, nil
}

// RevokeToken marca un refresh token como revocado
func (r *RefreshTokenRepository) RevokeToken(token string) error {
	result := r.db.Model(&models.RefreshToken{}).Where("token = ?", token).Update("is_revoked", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}
	return nil
}

// RevokeAllByUserID revoca todos los refresh tokens de un usuario
// DEPRECATED: Ya no se usa en la estrategia híbrida de logout.
// Se prefiere DeleteByUserID para eliminar físicamente y ahorrar espacio en BD
// para usuarios normales que no requieren auditoría de tokens.
//func (r *RefreshTokenRepository) RevokeAllByUserID(userID uint) error {
//	result := r.db.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("is_revoked", true)
//	return result.Error
//}

// DeleteByUserID elimina físicamente todos los refresh tokens de un usuario
// Usado en logout para limpiar tokens innecesarios de usuarios normales
func (r *RefreshTokenRepository) DeleteByUserID(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})
	return result.Error
}

// CleanupExpired elimina tokens expirados de la base de datos
func (r *RefreshTokenRepository) CleanupExpired() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{})
	return result.Error
}
