package models

import (
	"time"
)

// RefreshToken representa un token de actualización para mantener sesiones activas
type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"unique;not null;index"`
	UserID    uint      `json:"userId" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
	IsRevoked bool      `json:"isRevoked" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}