package services

import (
	"errors"
	"fmt"

	"backend-academi/internal/models"
	"backend-academi/internal/repository"
)

// QuestionnaireService interfaz para el servicio de cuestionarios
type QuestionnaireService interface {
	GetInitialQuestionnaire() (*models.QuestionnaireInfo, error)
	SubmitInitialQuestionnaire(userID uint, request *models.InitialQuestionnaireRequest) error
	HasUserCompletedInitial(userID uint) (bool, error)
	GetUserInitialResponse(userID uint) (*models.QuestionnaireResponse, error)
}

// questionnaireService implementación del servicio
type questionnaireService struct {
	questionnaireRepo repository.QuestionnaireRepository
	userRepo          repository.UserRepositoryInterface
}

// NewQuestionnaireService crea una nueva instancia del servicio
func NewQuestionnaireService(
	questionnaireRepo repository.QuestionnaireRepository,
	userRepo repository.UserRepositoryInterface,
) QuestionnaireService {
	return &questionnaireService{
		questionnaireRepo: questionnaireRepo,
		userRepo:          userRepo,
	}
}

// GetInitialQuestionnaire obtiene la información del cuestionario inicial
func (s *questionnaireService) GetInitialQuestionnaire() (*models.QuestionnaireInfo, error) {
	questionnaire, err := s.questionnaireRepo.GetByType(models.InitialQuestionnaireType)
	if err != nil {
		return nil, fmt.Errorf("error obtaining initial questionnaire: %w", err)
	}

	// Construir la información del cuestionario con las preguntas
	questionnaireInfo := &models.QuestionnaireInfo{
		ID:          questionnaire.ID,
		Type:        questionnaire.Type,
		Title:       questionnaire.Title,
		Description: questionnaire.Description,
		Questions:   s.buildInitialQuestionnaireQuestions(),
	}

	return questionnaireInfo, nil
}

// SubmitInitialQuestionnaire procesa y guarda las respuestas del cuestionario inicial
func (s *questionnaireService) SubmitInitialQuestionnaire(userID uint, request *models.InitialQuestionnaireRequest) error {
	// Verificar que el usuario existe
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Obtener el cuestionario inicial
	questionnaire, err := s.questionnaireRepo.GetByType(models.InitialQuestionnaireType)
	if err != nil {
		return fmt.Errorf("initial questionnaire not found: %w", err)
	}

	// Verificar si el usuario ya completó el cuestionario
	hasCompleted, err := s.questionnaireRepo.HasUserCompletedQuestionnaire(userID, models.InitialQuestionnaireType)
	if err != nil {
		return fmt.Errorf("error checking questionnaire completion: %w", err)
	}

	if hasCompleted {
		return errors.New("user has already completed the initial questionnaire")
	}

	// Validar la request
	if err := s.validateInitialQuestionnaireRequest(request); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Convertir materias de interés a JSON
	subjectsJSON, err := repository.ConvertSubjectsToJSON(request.SubjectsOfInterest)
	if err != nil {
		return fmt.Errorf("error converting subjects to JSON: %w", err)
	}

	// Crear la respuesta
	response := &models.QuestionnaireResponse{
		UserID:             userID,
		QuestionnaireID:    questionnaire.ID,
		StudyHoursPerDay:   request.StudyHoursPerDay,
		TimePreference:     request.TimePreference,
		PrimaryGoal:        request.PrimaryGoal,
		SecondaryGoal:      request.SecondaryGoal,
		CurrentLevel:       request.CurrentLevel,
		SubjectsOfInterest: subjectsJSON,
		AdditionalComments: request.AdditionalComments,
	}

	// Guardar en la base de datos
	if err := s.questionnaireRepo.CreateResponse(response); err != nil {
		return fmt.Errorf("error saving questionnaire response: %w", err)
	}

	// Actualizar el perfil del estudiante si existe
	if user.Role == models.StudentRole {
		if err := s.updateStudentProfile(userID, request); err != nil {
			// Log el error pero no fallar la operación principal
			fmt.Printf("Warning: Could not update student profile: %v\n", err)
		}
	}

	return nil
}

// HasUserCompletedInitial verifica si el usuario completó el cuestionario inicial
func (s *questionnaireService) HasUserCompletedInitial(userID uint) (bool, error) {
	return s.questionnaireRepo.HasUserCompletedQuestionnaire(userID, models.InitialQuestionnaireType)
}

// GetUserInitialResponse obtiene la respuesta del usuario al cuestionario inicial
func (s *questionnaireService) GetUserInitialResponse(userID uint) (*models.QuestionnaireResponse, error) {
	questionnaire, err := s.questionnaireRepo.GetByType(models.InitialQuestionnaireType)
	if err != nil {
		return nil, fmt.Errorf("initial questionnaire not found: %w", err)
	}

	response, err := s.questionnaireRepo.GetUserResponse(userID, questionnaire.ID)
	if err != nil {
		return nil, fmt.Errorf("user response not found: %w", err)
	}

	return response, nil
}

// buildInitialQuestionnaireQuestions construye las preguntas del cuestionario inicial
func (s *questionnaireService) buildInitialQuestionnaireQuestions() []models.QuestionInfo {
	return []models.QuestionInfo{
		{
			ID:       "study_hours_per_day",
			Text:     "¿Cuántas horas puedes dedicar al estudio por día?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "1", Label: "1 hora"},
				{Value: "2", Label: "2 horas"},
				{Value: "3", Label: "3 horas"},
				{Value: "4", Label: "4 horas"},
				{Value: "5", Label: "5 horas"},
				{Value: "6", Label: "6+ horas"},
			},
		},
		{
			ID:       "time_preference",
			Text:     "¿Cuál es tu horario preferido para estudiar?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "morning", Label: "Mañana (6:00 - 12:00)"},
				{Value: "afternoon", Label: "Tarde (12:00 - 18:00)"},
				{Value: "evening", Label: "Noche (18:00 - 22:00)"},
				{Value: "night", Label: "Madrugada (22:00 - 6:00)"},
			},
		},
		{
			ID:       "primary_goal",
			Text:     "¿Cuál es tu objetivo principal?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "pass_exam", Label: "Aprobar un examen específico"},
				{Value: "reinforce_course", Label: "Reforzar conocimientos de un curso"},
				{Value: "learn_new_topic", Label: "Aprender un tema nuevo"},
				{Value: "improve_grades", Label: "Mejorar mis calificaciones"},
				{Value: "prepare_university", Label: "Prepararme para la universidad"},
			},
		},
		{
			ID:       "secondary_goal",
			Text:     "¿Tienes un objetivo secundario? (Opcional)",
			Type:     "select",
			Required: false,
			Options: []models.OptionInfo{
				{Value: "pass_exam", Label: "Aprobar un examen específico"},
				{Value: "reinforce_course", Label: "Reforzar conocimientos de un curso"},
				{Value: "learn_new_topic", Label: "Aprender un tema nuevo"},
				{Value: "improve_grades", Label: "Mejorar mis calificaciones"},
				{Value: "prepare_university", Label: "Prepararme para la universidad"},
			},
		},
		{
			ID:       "current_level",
			Text:     "¿En qué nivel académico te encuentras? (Opcional)",
			Type:     "select",
			Required: false,
			Options: []models.OptionInfo{
				{Value: "primary", Label: "Primaria"},
				{Value: "secondary", Label: "Secundaria"},
				{Value: "high_school", Label: "Preparatoria"},
				{Value: "university", Label: "Universidad"},
				{Value: "other", Label: "Otro"},
			},
		},
		{
			ID:       "subjects_of_interest",
			Text:     "¿Qué materias te interesan más? (Opcional, puedes seleccionar varias)",
			Type:     "multiple",
			Required: false,
			Options: []models.OptionInfo{
				{Value: "mathematics", Label: "Matemáticas"},
				{Value: "physics", Label: "Física"},
				{Value: "chemistry", Label: "Química"},
				{Value: "biology", Label: "Biología"},
				{Value: "spanish", Label: "Español"},
				{Value: "english", Label: "Inglés"},
				{Value: "history", Label: "Historia"},
				{Value: "geography", Label: "Geografía"},
				{Value: "computer_science", Label: "Informática"},
				{Value: "arts", Label: "Artes"},
			},
		},
		{
			ID:       "additional_comments",
			Text:     "¿Hay algo más que te gustaría que sepamos? (Opcional)",
			Type:     "text",
			Required: false,
		},
	}
}

// validateInitialQuestionnaireRequest valida los datos del request
func (s *questionnaireService) validateInitialQuestionnaireRequest(request *models.InitialQuestionnaireRequest) error {
	if request.StudyHoursPerDay < 1 || request.StudyHoursPerDay > 12 {
		return errors.New("study hours per day must be between 1 and 12")
	}

	validTimePreferences := []models.StudyTimePreference{
		models.MorningPreference,
		models.AfternoonPreference,
		models.EveningPreference,
		models.NightPreference,
	}
	
	isValidTimePreference := false
	for _, valid := range validTimePreferences {
		if request.TimePreference == valid {
			isValidTimePreference = true
			break
		}
	}
	
	if !isValidTimePreference {
		return errors.New("invalid time preference")
	}

	validGoals := []models.StudyGoalType{
		models.PassExamGoal,
		models.ReinforceCourseGoal,
		models.LearnNewTopicGoal,
		models.ImproveGradesGoal,
		models.PrepareUniversityGoal,
	}
	
	isValidPrimaryGoal := false
	for _, valid := range validGoals {
		if request.PrimaryGoal == valid {
			isValidPrimaryGoal = true
			break
		}
	}
	
	if !isValidPrimaryGoal {
		return errors.New("invalid primary goal")
	}

	// Validar objetivo secundario si está presente
	if request.SecondaryGoal != nil {
		isValidSecondaryGoal := false
		for _, valid := range validGoals {
			if *request.SecondaryGoal == valid {
				isValidSecondaryGoal = true
				break
			}
		}
		
		if !isValidSecondaryGoal {
			return errors.New("invalid secondary goal")
		}
	}

	return nil
}

// updateStudentProfile actualiza el perfil del estudiante con información del cuestionario
func (s *questionnaireService) updateStudentProfile(userID uint, request *models.InitialQuestionnaireRequest) error {
	// Esta función se puede implementar más tarde cuando tengamos más lógica de estudiantes
	// Por ahora solo retornamos nil
	return nil
}