package models

import (
	"time"
)

// RefreshToken representa un token de actualización para mantener sesiones activas
type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"unique;not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}