package models

import (
	"time"
)

// TokenBlacklist representa tokens JWT invalidados para usuarios privilegiados (admin/teacher)
type TokenBlacklist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	JTI       string    `json:"jti" gorm:"unique;not null;index"` // JWT ID claim
	Token     string    `json:"token" gorm:"not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}