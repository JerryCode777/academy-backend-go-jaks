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
		UserID:                   userID,
		QuestionnaireID:          questionnaire.ID,
		StudyHoursPerDay:         request.StudyHoursPerDay,
		TimePreference:           request.TimePreference,
		PrimaryGoal:              request.PrimaryGoal,
		CurrentLevel:             request.CurrentLevel,
		SubjectsOfInterest:       subjectsJSON,
		ExamPreparationExperience: request.ExamPreparationExperience,
		AdditionalComments:       request.AdditionalComments,
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
			Text:     "¿Cuántas horas diarias puedes dedicar a prepararte para el examen de admisión?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "2", Label: "2 horas - Preparación básica"},
				{Value: "4", Label: "4 horas - Preparación estándar"},
				{Value: "6", Label: "6 horas - Preparación intensiva"},
				{Value: "8", Label: "8 horas - Preparación dedicada completa"},
			},
		},
		{
			ID:       "time_preference",
			Text:     "¿En qué momento del día prefieres estudiar para tu examen de admisión?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "morning", Label: "Mañana (6:00 - 12:00) - Máxima concentración"},
				{Value: "afternoon", Label: "Tarde (12:00 - 18:00) - Horario regular"},
				{Value: "evening", Label: "Noche (18:00 - 22:00) - Ambiente tranquilo"},
				{Value: "night", Label: "Madrugada (22:00 - 6:00) - Sin interrupciones"},
			},
		},
		{
			ID:       "primary_goal",
			Text:     "¿Cuál es tu objetivo principal con esta plataforma?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "unsa_admission", Label: "Ingresar a la UNSA - Universidad Nacional de San Agustín"},
				{Value: "other_national_university", Label: "Ingresar a otra universidad nacional"},
				{Value: "private_university", Label: "Ingresar a universidad privada"},
				{Value: "improve_exam_scores", Label: "Mejorar mis puntajes en simulacros"},
				{Value: "reinforce_knowledge", Label: "Reforzar conocimientos preuniversitarios"},
			},
		},
		{
			ID:       "current_level",
			Text:     "¿En qué año de estudios te encuentras actualmente?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "grade_4", Label: "4to de Secundaria - Preparándome con anticipación"},
				{Value: "grade_5", Label: "5to de Secundaria - Año de decisión"},
				{Value: "graduate", Label: "Egresado - Enfocado en el ingreso universitario"},
				{Value: "working_student", Label: "Estudiante-trabajador - Preparación flexible"},
			},
		},
		{
			ID:       "subjects_of_interest",
			Text:     "¿Qué carreras o áreas de estudio te interesan más? (selecciona todas las que apliquen)",
			Type:     "multiple",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "engineering", Label: "Ingenierías - Civil, Sistemas, Industrial, etc."},
				{Value: "medicine", Label: "Medicina - Medicina Humana, Enfermería"},
				{Value: "sciences", Label: "Ciencias Exactas - Matemáticas, Física, Química"},
				{Value: "economics", Label: "Ciencias Económicas - Economía, Administración"},
				{Value: "social_sciences", Label: "Ciencias Sociales - Derecho, Psicología, Educación"},
				{Value: "arts_humanities", Label: "Arte y Humanidades - Literatura, Historia, Filosofía"},
				{Value: "not_sure", Label: "Aún no estoy seguro - Necesito orientación"},
			},
		},
		{
			ID:       "exam_preparation_experience",
			Text:     "¿Cuánta experiencia tienes preparándote para exámenes de admisión?",
			Type:     "select",
			Required: true,
			Options: []models.OptionInfo{
				{Value: "beginner", Label: "Principiante - Es mi primera vez"},
				{Value: "some_experience", Label: "Algo de experiencia - He estudiado por mi cuenta"},
				{Value: "academy_experience", Label: "Con experiencia - He llevado academia preuniversitaria"},
				{Value: "multiple_attempts", Label: "Experimentado - He rendido exámenes anteriormente"},
			},
		},
		{
			ID:       "additional_comments",
			Text:     "¿Hay algo específico sobre tu preparación universitaria que te gustaría que sepamos?",
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
		// Nuevos objetivos preuniversitarios
		models.UNSAAdmissionGoal,
		models.OtherNationalUniversityGoal,
		models.PrivateUniversityGoal,
		models.ImproveExamScoresGoal,
		models.ReinforceKnowledgeGoal,
		// Compatibilidad con versiones anteriores
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