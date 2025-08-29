package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"backend-academi/internal/models"
)

// QuestionnaireRepository interfaz para operaciones del repositorio de cuestionarios
type QuestionnaireRepository interface {
	GetByType(questionnaireType models.QuestionnaireType) (*models.Questionnaire, error)
	CreateResponse(response *models.QuestionnaireResponse) error
	GetStudentResponse(studentID uint, questionnaireID uint) (*models.QuestionnaireResponse, error)
	HasStudentCompletedQuestionnaire(studentID uint, questionnaireType models.QuestionnaireType) (bool, error)
	CreateInitialQuestionnaire() error
}

// questionnaireRepository implementación del repositorio
type questionnaireRepository struct {
	db *gorm.DB
}

// NewQuestionnaireRepository crea una nueva instancia del repositorio
func NewQuestionnaireRepository(db *gorm.DB) QuestionnaireRepository {
	repo := &questionnaireRepository{db: db}
	
	// Crear el cuestionario inicial si no existe
	if err := repo.CreateInitialQuestionnaire(); err != nil {
		// Log el error pero no fallar la aplicación
		fmt.Printf("Warning: Could not create initial questionnaire: %v\n", err)
	}
	
	return repo
}

// GetByType obtiene un cuestionario por su tipo
func (r *questionnaireRepository) GetByType(questionnaireType models.QuestionnaireType) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	
	err := r.db.Where("type = ? AND is_active = ?", questionnaireType, true).First(&questionnaire).Error
	if err != nil {
		return nil, err
	}
	
	return &questionnaire, nil
}

// CreateResponse crea una nueva respuesta de cuestionario
func (r *questionnaireRepository) CreateResponse(response *models.QuestionnaireResponse) error {
	// Convertir el slice de materias a JSON para almacenar
	if response.SubjectsOfInterest == "" && len(response.SubjectsOfInterest) > 0 {
		// Si viene como array desde el request, convertir a JSON
		// Esto se manejará en el service
	}
	
	response.CompletedAt = time.Now()
	
	return r.db.Create(response).Error
}

// GetStudentResponse obtiene la respuesta de un estudiante para un cuestionario específico
func (r *questionnaireRepository) GetStudentResponse(studentID uint, questionnaireID uint) (*models.QuestionnaireResponse, error) {
	var response models.QuestionnaireResponse
	
	err := r.db.Where("student_id = ? AND questionnaire_id = ?", studentID, questionnaireID).
		Preload("Student").
		Preload("Questionnaire").
		First(&response).Error
		
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}

// HasStudentCompletedQuestionnaire verifica si un estudiante ha completado un tipo de cuestionario
func (r *questionnaireRepository) HasStudentCompletedQuestionnaire(studentID uint, questionnaireType models.QuestionnaireType) (bool, error) {
	var count int64
	
	err := r.db.Table("questionnaire_responses").
		Joins("JOIN questionnaires ON questionnaire_responses.questionnaire_id = questionnaires.id").
		Where("questionnaire_responses.student_id = ? AND questionnaires.type = ?", studentID, questionnaireType).
		Count(&count).Error
		
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// CreateInitialQuestionnaire crea el cuestionario inicial si no existe
func (r *questionnaireRepository) CreateInitialQuestionnaire() error {
	// Verificar si ya existe
	var count int64
	err := r.db.Model(&models.Questionnaire{}).Where("type = ?", models.InitialQuestionnaireType).Count(&count).Error
	if err != nil {
		return err
	}
	
	// Si ya existe, no crear otro
	if count > 0 {
		return nil
	}
	
	// Crear el cuestionario inicial
	initialQuestionnaire := &models.Questionnaire{
		Type:        models.InitialQuestionnaireType,
		Title:       "Cuestionario Inicial - Conoce tu perfil de estudio",
		Description: "Este cuestionario nos ayuda a personalizar tu experiencia de aprendizaje según tus necesidades y preferencias.",
		IsActive:    true,
	}
	
	return r.db.Create(initialQuestionnaire).Error
}

// ConvertSubjectsToJSON convierte un slice de strings a JSON
func ConvertSubjectsToJSON(subjects []string) (string, error) {
	if len(subjects) == 0 {
		return "[]", nil
	}
	
	jsonBytes, err := json.Marshal(subjects)
	if err != nil {
		return "", err
	}
	
	return string(jsonBytes), nil
}

// ConvertSubjectsFromJSON convierte JSON a slice de strings
func ConvertSubjectsFromJSON(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return []string{}, nil
	}
	
	var subjects []string
	err := json.Unmarshal([]byte(jsonStr), &subjects)
	if err != nil {
		return nil, err
	}
	
	return subjects, nil
}