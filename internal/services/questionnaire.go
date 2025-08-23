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
			Text:     "¿Cuánto tiempo puedes dedicar al estudio cada día?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "1", Label: "1 hora - Sesiones cortas y efectivas"},
				{Value: "2", Label: "2 horas - Balance perfecto"},
				{Value: "3", Label: "3 horas - Estudio profundo"},
				{Value: "4", Label: "4 horas - Preparación intensiva"},
				{Value: "5", Label: "5+ horas - Dedicación completa"},
			},
		},
		{
			ID:       "time_preference",
			Text:     "¿En qué momento del día te sientes más concentrado?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "morning", Label: "Mañana temprano - Mente fresca y energizada"},
				{Value: "afternoon", Label: "Tarde - Después del almuerzo"},
				{Value: "evening", Label: "Noche - Ambiente tranquilo"},
				{Value: "night", Label: "Madrugada - Sin distracciones"},
			},
		},
		{
			ID:       "primary_goal",
			Text:     "¿Qué te motiva a estudiar en este momento?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "pass_exam", Label: "Aprobar un examen importante"},
				{Value: "improve_grades", Label: "Mejorar mis calificaciones generales"},
				{Value: "learn_new_topic", Label: "Explorar nuevas materias"},
				{Value: "prepare_university", Label: "Prepararme para la universidad"},
				{Value: "reinforce_course", Label: "Reforzar lo que ya sé"},
			},
		},
		{
			ID:       "current_level",
			Text:     "¿En qué etapa de tu educación te encuentras?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "secondary", Label: "Secundaria - Construyendo bases sólidas"},
				{Value: "high_school", Label: "Preparatoria - Preparándome para el siguiente nivel"},
				{Value: "university", Label: "Universidad - Especializándome"},
				{Value: "other", Label: "Otro - Aprendizaje continuo"},
			},
		},
		{
			ID:       "subjects_of_interest",
			Text:     "¿Qué materias te llaman más la atención?",
			Type:     "multiple",
			Required: false,
			Options: []models.OptionInfo{
				{Value: "mathematics", Label: "Matemáticas - Lógica y resolución de problemas"},
				{Value: "sciences", Label: "Ciencias - Física, Química, Biología"},
				{Value: "languages", Label: "Idiomas - Español, Inglés, otros"},
				{Value: "social", Label: "Ciencias Sociales - Historia, Geografía"},
				{Value: "technology", Label: "Tecnología - Informática, Programación"},
				{Value: "arts", Label: "Artes - Creatividad y expresión"},
			},
		},
		{
			ID:       "additional_comments",
			Text:     "¿Tienes algún desafío específico que te gustaría abordar?",
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


	return nil
}

// updateStudentProfile actualiza el perfil del estudiante con información del cuestionario
func (s *questionnaireService) updateStudentProfile(userID uint, request *models.InitialQuestionnaireRequest) error {
	// Esta función se puede implementar más tarde cuando tengamos más lógica de estudiantes
	// Por ahora solo retornamos nil
	return nil
}