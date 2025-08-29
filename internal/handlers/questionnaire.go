package handlers

import (
	"encoding/json"
	"net/http"

	"backend-academi/internal/models"
	"backend-academi/internal/services"
	"backend-academi/pkg/utils"
)

// QuestionnaireHandler maneja las rutas relacionadas con cuestionarios
type QuestionnaireHandler struct {
	questionnaireService services.QuestionnaireService
}

// NewQuestionnaireHandler crea una nueva instancia del handler
func NewQuestionnaireHandler(questionnaireService services.QuestionnaireService) *QuestionnaireHandler {
	return &QuestionnaireHandler{
		questionnaireService: questionnaireService,
	}
}

// GetInitialQuestionnaire maneja GET /api/v1/questionnaire/initial
func (h *QuestionnaireHandler) GetInitialQuestionnaire(w http.ResponseWriter, r *http.Request) {
	questionnaire, err := h.questionnaireService.GetInitialQuestionnaire()
	if err != nil {
		http.Error(w, "Error obtaining questionnaire: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Initial questionnaire obtained successfully",
		"data":    questionnaire,
	})
}

// SubmitInitialQuestionnaire maneja POST /api/v1/questionnaire/initial/submit
func (h *QuestionnaireHandler) SubmitInitialQuestionnaire(w http.ResponseWriter, r *http.Request) {
	// Obtener las claims del contexto
	claims, ok := r.Context().Value("user_claims").(*utils.Claims)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	studentID := claims.UserID // UserID del JWT corresponde al StudentID

	// Decodificar el request
	var request models.InitialQuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Procesar la respuesta
	if err := h.questionnaireService.SubmitInitialQuestionnaire(studentID, &request); err != nil {
		// Diferentes códigos de error según el tipo
		if err.Error() == "student has already completed the initial questionnaire" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		
		if err.Error()[:15] == "invalid request" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		http.Error(w, "Error processing questionnaire: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Questionnaire submitted successfully",
	})
}

// CheckInitialCompletion maneja GET /api/v1/questionnaire/initial/status
func (h *QuestionnaireHandler) CheckInitialCompletion(w http.ResponseWriter, r *http.Request) {
	// Obtener las claims del contexto
	claims, ok := r.Context().Value("user_claims").(*utils.Claims)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	studentID := claims.UserID // UserID del JWT corresponde al StudentID

	hasCompleted, err := h.questionnaireService.HasStudentCompletedInitial(studentID)
	if err != nil {
		http.Error(w, "Error checking questionnaire status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"message":      "Questionnaire status obtained successfully",
		"has_completed": hasCompleted,
	})
}

// GetUserInitialResponse maneja GET /api/v1/questionnaire/initial/response
func (h *QuestionnaireHandler) GetUserInitialResponse(w http.ResponseWriter, r *http.Request) {
	// Obtener las claims del contexto
	claims, ok := r.Context().Value("user_claims").(*utils.Claims)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	
	studentID := claims.UserID // UserID del JWT corresponde al StudentID

	response, err := h.questionnaireService.GetStudentInitialResponse(studentID)
	if err != nil {
		http.Error(w, "Error obtaining user response: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User response obtained successfully",
		"data":    response,
	})
}

// GetInitialQuestionnairePublic maneja GET /api/v1/questionnaire/initial/public
// Esta ruta no requiere autenticación para que usuarios no logueados puedan ver el cuestionario
func (h *QuestionnaireHandler) GetInitialQuestionnairePublic(w http.ResponseWriter, r *http.Request) {
	questionnaire, err := h.questionnaireService.GetInitialQuestionnaire()
	if err != nil {
		http.Error(w, "Error obtaining questionnaire: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Initial questionnaire obtained successfully",
		"data":    questionnaire,
	})
}