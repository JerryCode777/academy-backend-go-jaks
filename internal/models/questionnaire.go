package models

import (
	"time"
)

// QuestionnaireType define los tipos de cuestionarios disponibles
type QuestionnaireType string

const (
	InitialQuestionnaireType QuestionnaireType = "initial"
	AssessmentQuestionnaireType QuestionnaireType = "assessment"
)

// StudyTimePreference representa las opciones de horario de estudio
type StudyTimePreference string

const (
	MorningPreference StudyTimePreference = "morning"
	AfternoonPreference StudyTimePreference = "afternoon"
	EveningPreference StudyTimePreference = "evening"
	NightPreference StudyTimePreference = "night"
)

// StudyGoalType representa los objetivos principales de estudio
type StudyGoalType string

const (
	PassExamGoal StudyGoalType = "pass_exam"
	ReinforceCourseGoal StudyGoalType = "reinforce_course"
	LearnNewTopicGoal StudyGoalType = "learn_new_topic"
	ImproveGradesGoal StudyGoalType = "improve_grades"
	PrepareUniversityGoal StudyGoalType = "prepare_university"
)

// Questionnaire representa un cuestionario del sistema
type Questionnaire struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Type        QuestionnaireType `json:"type" gorm:"not null"`
	Title       string            `json:"title" gorm:"not null"`
	Description string            `json:"description" gorm:"text"`
	IsActive    bool              `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// QuestionnaireResponse representa las respuestas de un usuario a un cuestionario
type QuestionnaireResponse struct {
	ID              uint                `json:"id" gorm:"primaryKey"`
	UserID          uint                `json:"user_id" gorm:"not null"`
	User            User                `json:"user" gorm:"foreignKey:UserID"`
	QuestionnaireID uint                `json:"questionnaire_id" gorm:"not null"`
	Questionnaire   Questionnaire       `json:"questionnaire" gorm:"foreignKey:QuestionnaireID"`
	
	// Respuestas específicas del cuestionario inicial
	StudyHoursPerDay    int                 `json:"study_hours_per_day" gorm:"not null"`
	TimePreference      StudyTimePreference `json:"time_preference" gorm:"not null"`
	PrimaryGoal         StudyGoalType       `json:"primary_goal" gorm:"not null"`
	SecondaryGoal       *StudyGoalType      `json:"secondary_goal,omitempty" gorm:"null"`
	
	// Campos adicionales
	CurrentLevel        string              `json:"current_level" gorm:"null"`        // Nivel académico actual
	SubjectsOfInterest  string              `json:"subjects_of_interest" gorm:"text"` // JSON array de materias
	AdditionalComments  string              `json:"additional_comments" gorm:"text"`  // Comentarios adicionales
	
	CompletedAt         time.Time           `json:"completed_at"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

// InitialQuestionnaireRequest estructura para recibir las respuestas del cuestionario inicial
type InitialQuestionnaireRequest struct {
	StudyHoursPerDay   int                 `json:"study_hours_per_day" validate:"required,min=1,max=12"`
	TimePreference     StudyTimePreference `json:"time_preference" validate:"required"`
	PrimaryGoal        StudyGoalType       `json:"primary_goal" validate:"required"`
	SecondaryGoal      *StudyGoalType      `json:"secondary_goal,omitempty"`
	CurrentLevel       string              `json:"current_level,omitempty"`
	SubjectsOfInterest []string            `json:"subjects_of_interest,omitempty"`
	AdditionalComments string              `json:"additional_comments,omitempty"`
}

// QuestionnaireInfo estructura para enviar información del cuestionario al frontend
type QuestionnaireInfo struct {
	ID          uint              `json:"id"`
	Type        QuestionnaireType `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Questions   []QuestionInfo    `json:"questions"`
}

// QuestionInfo estructura para representar una pregunta
type QuestionInfo struct {
	ID       string        `json:"id"`
	Text     string        `json:"text"`
	Type     string        `json:"type"` // "select", "number", "text", "multiple"
	Required bool          `json:"required"`
	Options  []OptionInfo  `json:"options,omitempty"`
}

// OptionInfo estructura para representar opciones de respuesta
type OptionInfo struct {
	Value string `json:"value"`
	Label string `json:"label"`
}