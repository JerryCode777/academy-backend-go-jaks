package models

import (
	"time"
)

// ============================================================================
// PREGUNTAS BASE (REUTILIZABLES)
// ============================================================================

// Question representa una pregunta base que pertenece a un topic (REUTILIZABLE)
type Question struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	TopicID      uint         `json:"topic_id" gorm:"not null"`
	Topic        Topic        `json:"topic" gorm:"foreignKey:TopicID"`
	QuestionText string       `json:"question_text" gorm:"text;not null"`
	QuestionType string       `json:"question_type" gorm:"not null"` // multiple_choice, true_false, open
	Difficulty   int          `json:"difficulty" gorm:"check:difficulty >= 1 AND difficulty <= 10"`
	Points       int          `json:"points" gorm:"default:1"`
	Explanation  string       `json:"explanation" gorm:"text"` // Explicación de la respuesta
	OrderInTopic *int         `json:"order_in_topic,omitempty"` // Orden sugerido (no obligatorio)
	Tags         string       `json:"tags" gorm:"text"` // JSON array: ["algebra", "equations"]
	
	// Campos adicionales para importación
	TemporalID   string       `json:"temporal_id,omitempty" gorm:"index"` // id_temporal del JSON
	ImageURL     string       `json:"image_url,omitempty"`               // imagen del JSON
	
	IsActive     bool         `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	
	// Relaciones
	Options []QuestionOption `json:"options" gorm:"foreignKey:QuestionID"`
}

// QuestionOption representa una opción de respuesta para preguntas
type QuestionOption struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	QuestionID   uint      `json:"question_id" gorm:"not null"`
	Question     Question  `json:"question" gorm:"foreignKey:QuestionID"`
	OptionText   string    `json:"option_text" gorm:"not null"`
	IsCorrect    bool      `json:"is_correct" gorm:"default:false"`
	Explanation  string    `json:"explanation" gorm:"text"` // Por qué es correcta/incorrecta
	OrderPosition int      `json:"order_position"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ============================================================================
// SISTEMA DE QUIZZES DINÁMICOS
// ============================================================================

// QuizTemplate define plantillas reutilizables para generar quizzes
type QuizTemplate struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Name               string    `json:"name" gorm:"not null"`
	Description        string    `json:"description" gorm:"text"`
	Type               string    `json:"type" gorm:"not null;check:type IN ('quick', 'simulation', 'practice', 'formal')"`
	CourseID           *uint     `json:"course_id,omitempty"` // NULL = cross-course
	Course             *Course   `json:"course,omitempty" gorm:"foreignKey:CourseID"`
	TopicID            *uint     `json:"topic_id,omitempty"`  // NULL = cross-topic
	Topic              *Topic    `json:"topic,omitempty" gorm:"foreignKey:TopicID"`
	MaxQuestions       int       `json:"max_questions" gorm:"not null"`
	DefaultTimeMinutes int       `json:"default_time_minutes"`
	AllowCustomTime    bool      `json:"allow_custom_time" gorm:"default:true"`
	DifficultyRange    string    `json:"difficulty_range"` // "1-5", "6-10", "any"
	QuestionSelection  string    `json:"question_selection" gorm:"not null;check:question_selection IN ('random', 'difficulty_progressive', 'recent_mistakes', 'weakest_topics')"`
	IsActive           bool      `json:"is_active" gorm:"default:true"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// QuizInstance representa una instancia específica de quiz generada para un estudiante
type QuizInstance struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	StudentID      uint       `json:"student_id" gorm:"not null"`
	Student        Student    `json:"student" gorm:"foreignKey:StudentID"`
	TemplateID     uint       `json:"template_id" gorm:"not null"`
	Template       QuizTemplate `json:"template" gorm:"foreignKey:TemplateID"`
	Title          string     `json:"title"` // Personalizable por estudiante
	TimeLimit      int        `json:"time_limit"` // Minutos elegidos por usuario
	TotalQuestions int        `json:"total_questions" gorm:"not null"`
	Status         string     `json:"status" gorm:"not null;default:'generated';check:status IN ('generated', 'in_progress', 'completed', 'expired', 'abandoned')"`
	GeneratedAt    time.Time  `json:"generated_at" gorm:"not null"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	ExpiresAt      time.Time  `json:"expires_at" gorm:"not null"` // 24h después de generar
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	
	// Relaciones
	Questions []QuizInstanceQuestion `json:"questions" gorm:"foreignKey:QuizInstanceID"`
	Attempts  []QuizAttempt          `json:"attempts" gorm:"foreignKey:QuizInstanceID"`
}

// QuizInstanceQuestion tabla intermedia Many-to-Many con orden específico
type QuizInstanceQuestion struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	QuizInstanceID uint      `json:"quiz_instance_id" gorm:"not null"`
	QuizInstance   QuizInstance `json:"quiz_instance" gorm:"foreignKey:QuizInstanceID"`
	QuestionID     uint      `json:"question_id" gorm:"not null"` // REUTILIZA questions existentes
	Question       Question  `json:"question" gorm:"foreignKey:QuestionID"`
	OrderInQuiz    int       `json:"order_in_quiz" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
}

// QuizAttempt representa un intento de completar un quiz
type QuizAttempt struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	QuizInstanceID   uint       `json:"quiz_instance_id" gorm:"not null"`
	QuizInstance     QuizInstance `json:"quiz_instance" gorm:"foreignKey:QuizInstanceID"`
	StudentID        uint       `json:"student_id" gorm:"not null"`
	Student          Student    `json:"student" gorm:"foreignKey:StudentID"`
	AttemptNumber    int        `json:"attempt_number" gorm:"not null"`
	StartedAt        time.Time  `json:"started_at" gorm:"not null"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	TimeUsedMinutes  *int       `json:"time_used_minutes,omitempty"`
	Score            *float64   `json:"score,omitempty"` // 0.00-100.00
	QuestionsCorrect *int       `json:"questions_correct,omitempty"`
	QuestionsTotal   int        `json:"questions_total" gorm:"not null"`
	Status           string     `json:"status" gorm:"not null;default:'in_progress';check:status IN ('in_progress', 'completed', 'abandoned', 'expired')"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	
	// Relaciones
	Responses []QuizResponse `json:"responses" gorm:"foreignKey:QuizAttemptID"`
}

// QuizResponse representa una respuesta específica a una pregunta en un quiz
type QuizResponse struct {
	ID                uint             `json:"id" gorm:"primaryKey"`
	QuizAttemptID     uint             `json:"quiz_attempt_id" gorm:"not null"`
	QuizAttempt       QuizAttempt      `json:"quiz_attempt" gorm:"foreignKey:QuizAttemptID"`
	QuestionID        uint             `json:"question_id" gorm:"not null"` // REUTILIZA questions
	Question          Question         `json:"question" gorm:"foreignKey:QuestionID"`
	SelectedOptionID  *uint            `json:"selected_option_id,omitempty"` // NULL para preguntas abiertas
	SelectedOption    *QuestionOption  `json:"selected_option,omitempty" gorm:"foreignKey:SelectedOptionID"`
	OpenResponse      string           `json:"open_response" gorm:"text"` // Para preguntas abiertas
	IsCorrect         *bool            `json:"is_correct,omitempty"`
	TimeSpentSeconds  *int             `json:"time_spent_seconds,omitempty"`
	RespondedAt       time.Time        `json:"responded_at" gorm:"not null"`
	CreatedAt         time.Time        `json:"created_at"`
}

// ============================================================================
// ENUMS Y TIPOS
// ============================================================================

// QuizTemplateType define los tipos de plantillas de quiz
type QuizTemplateType string

const (
	QuickQuizType      QuizTemplateType = "quick"      // 5-10 preguntas rápidas
	SimulationQuizType QuizTemplateType = "simulation" // 60-90 preguntas tipo examen
	PracticeQuizType   QuizTemplateType = "practice"   // Práctica general
	FormalQuizType     QuizTemplateType = "formal"     // Evaluación formal
)

// QuizInstanceStatus define los estados de una instancia de quiz
type QuizInstanceStatus string

const (
	GeneratedStatus  QuizInstanceStatus = "generated"   // Recién generado
	InProgressStatus QuizInstanceStatus = "in_progress" // En progreso
	CompletedStatus  QuizInstanceStatus = "completed"   // Completado
	ExpiredStatus    QuizInstanceStatus = "expired"     // Expirado
	AbandonedStatus  QuizInstanceStatus = "abandoned"   // Abandonado
)

// QuestionSelectionType define cómo se seleccionan las preguntas
type QuestionSelectionType string

const (
	RandomSelection            QuestionSelectionType = "random"
	DifficultyProgressiveSelection QuestionSelectionType = "difficulty_progressive"
	RecentMistakesSelection    QuestionSelectionType = "recent_mistakes"
	WeakestTopicsSelection     QuestionSelectionType = "weakest_topics"
)

// ============================================================================
// REQUESTS Y RESPONSES PARA API
// ============================================================================

// CreateQuizInstanceRequest para crear una nueva instancia de quiz
type CreateQuizInstanceRequest struct {
	TemplateID       uint   `json:"template_id" validate:"required"`
	Title            string `json:"title,omitempty"`
	CustomTimeLimit  *int   `json:"custom_time_limit,omitempty"` // Minutos
	MaxQuestions     *int   `json:"max_questions,omitempty"`     // Override del template
}

// QuizInstanceInfo información resumida de una instancia
type QuizInstanceInfo struct {
	ID             uint                   `json:"id"`
	Title          string                 `json:"title"`
	Type           QuizTemplateType       `json:"type"`
	TotalQuestions int                    `json:"total_questions"`
	TimeLimit      int                    `json:"time_limit"`
	Status         QuizInstanceStatus     `json:"status"`
	GeneratedAt    time.Time              `json:"generated_at"`
	ExpiresAt      time.Time              `json:"expires_at"`
	Questions      []QuizQuestionInfo     `json:"questions,omitempty"`
}

// QuizQuestionInfo información de pregunta para el frontend
type QuizQuestionInfo struct {
	ID           uint                    `json:"id"`
	QuestionText string                  `json:"question_text"`
	QuestionType string                  `json:"question_type"`
	Points       int                     `json:"points"`
	OrderInQuiz  int                     `json:"order_in_quiz"`
	Options      []QuizQuestionOptionInfo `json:"options,omitempty"`
}

// QuizQuestionOptionInfo información de opción para el frontend
type QuizQuestionOptionInfo struct {
	ID           uint   `json:"id"`
	OptionText   string `json:"option_text"`
	OrderPosition int   `json:"order_position"`
}

// SubmitQuizResponseRequest para enviar respuesta a una pregunta
type SubmitQuizResponseRequest struct {
	QuestionID       uint   `json:"question_id" validate:"required"`
	SelectedOptionID *uint  `json:"selected_option_id,omitempty"`
	OpenResponse     string `json:"open_response,omitempty"`
	TimeSpentSeconds int    `json:"time_spent_seconds,omitempty"`
}

// QuizAttemptSummary resumen de un intento de quiz
type QuizAttemptSummary struct {
	ID               uint      `json:"id"`
	AttemptNumber    int       `json:"attempt_number"`
	StartedAt        time.Time `json:"started_at"`
	CompletedAt      *time.Time `json:"completed_at"`
	TimeUsedMinutes  *int      `json:"time_used_minutes"`
	Score            *float64  `json:"score"`
	QuestionsCorrect *int      `json:"questions_correct"`
	QuestionsTotal   int       `json:"questions_total"`
	Status           string    `json:"status"`
}