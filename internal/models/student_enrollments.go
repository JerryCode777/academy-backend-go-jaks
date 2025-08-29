package models

import (
	"time"
)

type StudentEnrollment struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	StudentID      uint       `json:"student_id" gorm:"not null"`
	Student        Student    `json:"student" gorm:"foreignKey:StudentID"`
	CourseID       uint       `json:"course_id" gorm:"not null"`
	Course         Course     `json:"course" gorm:"foreignKey:CourseID"`
	CurrentTopicID uint       `json:"current_topic_id,omitempty"`
	Topic          Topic      `json:"topic,omitempty" gorm:"foreignKey:CurrentTopicID"`
	EnrolledAt     time.Time  `json:"enrolled_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Status         string     `json:"status" gorm:"check:status IN ('in_progress','completed','dropped');default:'in_progress'"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	LastActivity   time.Time  `json:"last_activity"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// StudentTopicProgress trackea progreso LIGERO por topic (NO granular por pregunta)
type StudentTopicProgress struct {
	ID                   uint       `json:"id" gorm:"primaryKey"`
	StudentID            uint       `json:"student_id" gorm:"not null"`
	Student              Student    `json:"student" gorm:"foreignKey:StudentID"`
	TopicID              uint       `json:"topic_id" gorm:"not null"`
	Topic                Topic      `json:"topic" gorm:"foreignKey:TopicID"`
	
	// Solo contadores y estado general
	QuestionsAnswered    int        `json:"questions_answered" gorm:"default:0"`
	CorrectAnswers       int        `json:"correct_answers" gorm:"default:0"`
	LastQuestionID       *uint      `json:"last_question_id,omitempty"` // Para continuar donde quedó
	LastQuestion         *Question  `json:"last_question,omitempty" gorm:"foreignKey:LastQuestionID"`
	CompletionPercentage float64    `json:"completion_percentage" gorm:"default:0"`
	IsCompleted          bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	LastActivity         time.Time  `json:"last_activity" gorm:"not null"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}
