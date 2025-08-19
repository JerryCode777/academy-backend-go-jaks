package models

import (
	"time"
)

// Course representa un curso o materia
type Course struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"text"`
	Subject     string    `json:"subject" gorm:"not null"` // Matemáticas, Ciencias, etc.
	Difficulty  int       `json:"difficulty"`              // 1-10
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relaciones
	Topics []Topic `json:"topics" gorm:"foreignKey:CourseID"`
}

// Topic representa un tema dentro de un curso
type Topic struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CourseID    uint      `json:"course_id" gorm:"not null"`
	Course      Course    `json:"course" gorm:"foreignKey:CourseID"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"text"`
	Content     string    `json:"content" gorm:"text"`
	Order       int       `json:"order"`                   // Orden del tema en el curso
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}