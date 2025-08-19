package models

import (
	"time"
)

// User representa un usuario del sistema
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // No se serializa en JSON
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Role      UserRole  `json:"role" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRole define los roles de usuario disponibles
type UserRole string

const (
	StudentRole UserRole = "student"
	AdminRole   UserRole = "admin"
	TeacherRole UserRole = "teacher"
)

// Student representa los datos específicos de un estudiante
type Student struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Grade        int       `json:"grade"`                    // Grado académico
	Institution  string    `json:"institution"`              // Institución educativa
	TargetCareer string    `json:"target_career"`            // Carrera objetivo
	StudyGoals   string    `json:"study_goals" gorm:"text"`  // Metas de estudio
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}