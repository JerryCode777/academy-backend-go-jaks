package repository

import (
	"time"
	"backend-academi/internal/models"
	"gorm.io/gorm"
)

// TokenBlacklistRepository maneja todas las operaciones de base de datos para tokens en blacklist
type TokenBlacklistRepository struct {
	db *gorm.DB
}

// NewTokenBlacklistRepository crea una nueva instancia del repositorio de blacklist
func NewTokenBlacklistRepository(db *gorm.DB) *TokenBlacklistRepository {
	return &TokenBlacklistRepository{db: db}
}

// Add añade un token a la blacklist
func (r *TokenBlacklistRepository) Add(blacklistToken *models.TokenBlacklist) error {
	result := r.db.Create(blacklistToken)
	return result.Error
}

// IsBlacklisted verifica si un JTI está en la blacklist y no ha expirado
func (r *TokenBlacklistRepository) IsBlacklisted(jti string) (bool, error) {
	var count int64
	err := r.db.Model(&models.TokenBlacklist{}).Where("jti = ? AND expires_at > ?", jti, time.Now()).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CleanupExpired elimina tokens expirados de la blacklist
func (r *TokenBlacklistRepository) CleanupExpired() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&models.TokenBlacklist{})
	return result.Error
}