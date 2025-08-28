# 🎓 API Cuestionario Preuniversitario - Backend (27/08/2025)

## 📋 Resumen de Cambios Backend

### 🔄 **Actualización del Cuestionario Inicial**

#### **Objetivo Principal**
Transformar el cuestionario genérico en uno específicamente diseñado para estudiantes que se preparan para exámenes de admisión universitaria, con énfasis en la **UNSA (Universidad Nacional de San Agustín)**.

---

## 🏗️ **Cambios en la Estructura de Datos**

### **1. Nuevas Constantes de Objetivos**

**Archivo**: `/internal/models/questionnaire.go`

```go
// Nuevos objetivos preuniversitarios
const (
    UNSAAdmissionGoal           StudyGoalType = "unsa_admission"
    OtherNationalUniversityGoal StudyGoalType = "other_national_university"
    PrivateUniversityGoal       StudyGoalType = "private_university"
    ImproveExamScoresGoal       StudyGoalType = "improve_exam_scores"
    ReinforceKnowledgeGoal      StudyGoalType = "reinforce_knowledge"
)
```

### **2. Nuevo Campo en Modelo**

```go
type QuestionnaireResponse struct {
    // ... campos existentes
    ExamPreparationExperience string `json:"examPreparationExperience" gorm:"not null"`
    // ... resto de campos
}

type InitialQuestionnaireRequest struct {
    // ... campos existentes  
    ExamPreparationExperience string `json:"examPreparationExperience" validate:"required"`
    // ... resto de campos
}
```

---

## 📝 **Preguntas Actualizadas del Cuestionario**

### **Pregunta 1: Horas de Estudio**
```go
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
}
```

### **Pregunta 2: Horario Preferido**
```go
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
}
```

### **Pregunta 3: Objetivo Principal** ⭐
```go
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
}
```

### **Pregunta 4: Nivel Académico**
```go
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
}
```

### **Pregunta 5: Carreras de Interés** 🎓
```go
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
}
```

### **Pregunta 6: Experiencia** 🆕
```go
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
}
```

### **Pregunta 7: Comentarios**
```go
{
    ID:       "additional_comments",
    Text:     "¿Hay algo específico sobre tu preparación universitaria que te gustaría que sepamos?",
    Type:     "text",
    Required: false,
}
```

---

## 🔍 **Validaciones Actualizadas**

### **Objetivos Válidos**
```go
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
```

---

## 🌐 **Endpoints de API**

### **1. Obtener Cuestionario Inicial**
```http
GET /api/v1/questionnaire/initial
Authorization: Bearer <token>
```

**Respuesta:**
```json
{
    "status": "success",
    "data": {
        "id": 1,
        "type": "initial",
        "title": "Cuestionario de Preparación Preuniversitaria",
        "description": "Personaliza tu plan de estudios para el examen de admisión",
        "questions": [...]
    }
}
```

### **2. Enviar Respuestas**
```http
POST /api/v1/questionnaire/initial/submit
Authorization: Bearer <token>
Content-Type: application/json
```

**Body:**
```json
{
    "studyHoursPerDay": 4,
    "timePreference": "morning",
    "primaryGoal": "unsa_admission",
    "currentLevel": "grade_5", 
    "subjectsOfInterest": ["engineering", "sciences"],
    "examPreparationExperience": "some_experience",
    "additionalComments": "Necesito ayuda especialmente con matemáticas"
}
```

### **3. Verificar Estado**
```http
GET /api/v1/questionnaire/initial/status
Authorization: Bearer <token>
```

**Respuesta:**
```json
{
    "status": "success", 
    "data": {
        "hasCompleted": true
    }
}
```

### **4. Obtener Respuesta del Usuario**
```http
GET /api/v1/questionnaire/initial/response
Authorization: Bearer <token>
```

**Respuesta:**
```json
{
    "status": "success",
    "data": {
        "user": {
            "id": 1,
            "firstName": "Juan",
            "lastName": "Pérez",
            "email": "juan@example.com"
        },
        "questionnaire": {
            "id": 1,
            "type": "initial",
            "title": "Cuestionario de Preparación Preuniversitaria"
        },
        "studyHoursPerDay": 4,
        "timePreference": "morning",
        "primaryGoal": "unsa_admission",
        "currentLevel": "grade_5",
        "subjectsOfInterest": "[\"engineering\",\"sciences\"]",
        "examPreparationExperience": "some_experience",
        "additionalComments": "Necesito ayuda especialmente con matemáticas",
        "completedAt": "2025-08-27T08:00:00Z"
    }
}
```

---

## 🛠️ **Archivos Modificados**

### **Modelos**
- `internal/models/questionnaire.go`
  - ✅ Nuevas constantes `StudyGoalType` 
  - ✅ Campo `ExamPreparationExperience` agregado
  - ✅ Request struct actualizado

### **Servicios**  
- `internal/services/questionnaire.go`
  - ✅ Función `buildInitialQuestionnaireQuestions()` actualizada
  - ✅ Validaciones expandidas
  - ✅ Nuevo campo incluido en creación de respuesta

### **Base de Datos**
- **Migración necesaria** para agregar columna:
```sql
ALTER TABLE questionnaire_responses 
ADD COLUMN exam_preparation_experience VARCHAR(50) NOT NULL DEFAULT 'beginner';
```

---

## 📊 **Compatibilidad**

### **Retrocompatibilidad**
- ✅ Objetivos antiguos aún válidos
- ✅ Frontend existente funcional
- ✅ Migraciones automáticas de GORM

### **Nuevas Funcionalidades**
- ✅ Enfoque específico en preparación universitaria
- ✅ Opciones dirigidas a estudiantes peruanos
- ✅ Más contexto sobre experiencia previa

---

## 🧪 **Testing Manual**

### **Comandos de Prueba**
```bash
# 1. Obtener cuestionario
curl -X GET http://localhost:8080/api/v1/questionnaire/initial \
  -H "Authorization: Bearer <token>"

# 2. Enviar respuestas
curl -X POST http://localhost:8080/api/v1/questionnaire/initial/submit \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "studyHoursPerDay": 4,
    "timePreference": "morning", 
    "primaryGoal": "unsa_admission",
    "currentLevel": "grade_5",
    "subjectsOfInterest": ["engineering"],
    "examPreparationExperience": "beginner",
    "additionalComments": "Primera vez preparándome"
  }'

# 3. Verificar estado
curl -X GET http://localhost:8080/api/v1/questionnaire/initial/status \
  -H "Authorization: Bearer <token>"
```

---

## 🚀 **Estado del Desarrollo**

### **✅ Completado**
- [x] Preguntas rediseñadas para contexto preuniversitario
- [x] Nuevas constantes y validaciones
- [x] Campo de experiencia agregado
- [x] Compatibilidad mantenida
- [x] Documentación API actualizada

### **⏳ Pendiente**
- [ ] Testing completo de endpoints
- [ ] Verificación de migración de BD
- [ ] Validación del flujo completo

---

**Fecha**: 27 de Agosto, 2025  
**Versión**: v1.2.0 - Cuestionario Preuniversitario  
**Enfoque**: Preparación para UNSA y universidades peruanas