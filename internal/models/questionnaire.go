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
	UNSAAdmissionGoal StudyGoalType = "unsa_admission"
	OtherNationalUniversityGoal StudyGoalType = "other_national_university"
	PrivateUniversityGoal StudyGoalType = "private_university"
	ImproveExamScoresGoal StudyGoalType = "improve_exam_scores"
	ReinforceKnowledgeGoal StudyGoalType = "reinforce_knowledge"
	
	// Mantener compatibilidad con versiones anteriores
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
	IsActive    bool              `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// QuestionnaireResponse representa las respuestas de un usuario a un cuestionario
type QuestionnaireResponse struct {
	ID              uint                `json:"id" gorm:"primaryKey"`
	UserID          uint                `json:"userId" gorm:"not null"`
	User            User                `json:"user" gorm:"foreignKey:UserID"`
	QuestionnaireID uint                `json:"questionnaireId" gorm:"not null"`
	Questionnaire   Questionnaire       `json:"questionnaire" gorm:"foreignKey:QuestionnaireID"`
	
	// Respuestas específicas del cuestionario inicial
	StudyHoursPerDay    int                 `json:"studyHoursPerDay" gorm:"not null"`
	TimePreference      StudyTimePreference `json:"timePreference" gorm:"not null"`
	PrimaryGoal         StudyGoalType       `json:"primaryGoal" gorm:"not null"`
	
	// Campos adicionales
	CurrentLevel             string              `json:"currentLevel" gorm:"not null"`    // Nivel académico actual
	SubjectsOfInterest       string              `json:"subjectsOfInterest" gorm:"text"` // JSON array de materias
	ExamPreparationExperience string             `json:"examPreparationExperience" gorm:"not null"` // Experiencia en preparación
	AdditionalComments       string              `json:"additionalComments" gorm:"text"`  // Comentarios adicionales
	
	CompletedAt         time.Time           `json:"completedAt"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
}

// InitialQuestionnaireRequest estructura para recibir las respuestas del cuestionario inicial
type InitialQuestionnaireRequest struct {
	StudyHoursPerDay          int                 `json:"studyHoursPerDay" validate:"required,min=1,max=12"`
	TimePreference            StudyTimePreference `json:"timePreference" validate:"required"`
	PrimaryGoal               StudyGoalType       `json:"primaryGoal" validate:"required"`
	CurrentLevel              string              `json:"currentLevel" validate:"required"`
	SubjectsOfInterest        []string            `json:"subjectsOfInterest,omitempty"`
	ExamPreparationExperience string              `json:"examPreparationExperience" validate:"required"`
	AdditionalComments        string              `json:"additionalComments,omitempty"`
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