package models

import (
	"time"
)

// Quiz representa una evaluación o examen
type Quiz struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CourseID    uint      `json:"courseId" gorm:"not null"`
	Course      Course    `json:"course" gorm:"foreignKey:CourseID"`
	TopicID     *uint     `json:"topicId,omitempty"`
	Topic       *Topic    `json:"topic,omitempty" gorm:"foreignKey:TopicID"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"text"`
	TimeLimit   int       `json:"timeLimit"`              // En minutos
	Difficulty  int       `json:"difficulty"`              // 1-10
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	
	// Relaciones
	Questions []Question `json:"questions" gorm:"foreignKey:QuizID"`
}

// Question representa una pregunta en un quiz
type Question struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	QuizID       uint      `json:"quizId" gorm:"not null"`
	Quiz         Quiz      `json:"quiz" gorm:"foreignKey:QuizID"`
	QuestionText string    `json:"questionText" gorm:"text;not null"`
	QuestionType string    `json:"questionType" gorm:"not null"` // multiple_choice, true_false, open
	Points       int       `json:"points" gorm:"default:1"`
	Order        int       `json:"order"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	
	// Relaciones
	Options []QuestionOption `json:"options" gorm:"foreignKey:QuestionID"`
}

// QuestionOption representa una opción de respuesta para preguntas de opción múltiple
type QuestionOption struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	QuestionID uint      `json:"questionId" gorm:"not null"`
	Question   Question  `json:"question" gorm:"foreignKey:QuestionID"`
	OptionText string    `json:"optionText" gorm:"not null"`
	IsCorrect  bool      `json:"isCorrect" gorm:"default:false"`
	Order      int       `json:"order"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}