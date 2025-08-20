package auth

import (
	"errors"
	"backend-academi/internal/models"
	"backend-academi/internal/repository"
	"backend-academi/pkg/utils"
)

// AuthService maneja toda la lógica de autenticación y autorización
type AuthService struct {
	userRepo    *repository.UserRepository
	hashService *utils.HashService
	jwtService  *utils.JWTService
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(userRepo *repository.UserRepository, hashService *utils.HashService, jwtService *utils.JWTService) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		hashService: hashService,
		jwtService:  jwtService,
	}
}

// Register registra un nuevo usuario en el sistema
// Valida que el email no exista, hashea la contraseña y crea el usuario
func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	// 1. Verificar que el email no exista
	exists, err := s.userRepo.Exists(req.Email)
	if err != nil {
		return nil, errors.New("error checking email availability")
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// 2. Hashear contraseña
	hashedPassword, err := s.hashService.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error processing password")
	}

	// 3. Crear usuario con role por defecto de estudiante
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.StudentRole, // Role por defecto
		IsActive:  true,
	}

	// 4. Guardar en base de datos
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("error creating user")
	}

	// No retornar la contraseña hasheada
	user.Password = ""
	return user, nil
}

// Login autentica un usuario y genera un token JWT
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// 1. Buscar usuario por email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials") // No revelar si el email existe
	}

	// 2. Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// 3. Verificar contraseña
	if !s.hashService.VerifyPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// 4. Generar token JWT
	token, expiresAt, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.New("error generating token")
	}

	// 5. Preparar respuesta (sin contraseña)
	user.Password = ""
	response := &LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}

	return response, nil
}

// ValidateToken valida un token JWT y retorna los datos actuales del usuario
// Usado por middleware para verificar autenticación
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// 1. Validar token y extraer claims
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// 2. Obtener datos actuales del usuario desde BD
	// Esto asegura que si el usuario fue desactivado, no pueda usar tokens viejos
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 3. Verificar que el usuario siga activo
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// No retornar contraseña
	user.Password = ""
	return user, nil
}

// TODO: Método Logout a implementar después según decisión
// func (s *AuthService) Logout(tokenString string) error {
//     // Implementar blacklist si se decide añadir
//     return nil
// }