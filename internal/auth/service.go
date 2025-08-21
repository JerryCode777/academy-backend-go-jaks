package auth

import (
	"errors"
	"time"
	"crypto/rand"
	"encoding/hex"
	
	"backend-academi/internal/models"
	"backend-academi/internal/repository"
	"backend-academi/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthService maneja la lógica de autenticación y autorización
type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	blacklistRepo    *repository.TokenBlacklistRepository
	jwtService       *utils.JWTService
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	blacklistRepo *repository.TokenBlacklistRepository,
	jwtService *utils.JWTService,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		blacklistRepo:    blacklistRepo,
		jwtService:       jwtService,
	}
}

// Login autentica un usuario y genera un token JWT
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// 1. Buscar usuario por email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Verificar contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// 4. Generar JWT token (15 minutos)
	accessToken, expiresAt, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, errors.New("error generating access token")
	}

	// 5. Generar refresh token (7 días)
	refreshTokenString, err := s.generateRefreshToken()
	if err != nil {
		return nil, errors.New("error generating refresh token")
	}

	// 6. Guardar refresh token en BD
	refreshToken := &models.RefreshToken{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 días
		IsRevoked: false,
	}
	err = s.refreshTokenRepo.Create(refreshToken)
	if err != nil {
		return nil, errors.New("error saving refresh token")
	}

	// No retornar la contraseña en la respuesta
	user.Password = ""

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		User:         *user,
		ExpiresAt:    expiresAt,
	}, nil
}

// generateRefreshToken genera un token de refresh aleatorio
func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Register registra un nuevo usuario en el sistema
func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	// 1. Verificar que el email no exista
	exists, err := s.userRepo.Exists(req.Email)
	if err != nil {
		return nil, errors.New("error checking email existence")
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// 2. Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	// 3. Crear usuario (por defecto como student)
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.StudentRole,
		IsActive:  true,
	}

	// 4. Guardar en base de datos
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("error creating user")
	}

	// No retornar la contraseña
	user.Password = ""
	return user, nil
}

// ValidateToken valida un token JWT
func (s *AuthService) ValidateToken(token string) (*models.User, error) {
	// 1. Validar token JWT
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// 2. Buscar usuario por ID
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 3. Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// No retornar la contraseña
	user.Password = ""
	return user, nil
}

// RefreshToken actualiza un token de acceso usando un refresh token válido
func (s *AuthService) RefreshToken(req RefreshTokenRequest) (*LoginResponse, error) {
	// 1. Validar que el refresh token exista y no esté revocado/expirado
	refreshToken, err := s.refreshTokenRepo.GetByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// 2. Verificar que el usuario asociado esté activo
	if !refreshToken.User.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// 3. Generar nuevo access token
	accessToken, expiresAt, err := s.jwtService.GenerateToken(&refreshToken.User)
	if err != nil {
		return nil, errors.New("error generating access token")
	}

	// 4. Opcionalmente rotar el refresh token para mayor seguridad
	// Por ahora reutilizamos el mismo refresh token
	// En producción podrías implementar rotación de refresh tokens

	// No retornar la contraseña en la respuesta
	refreshToken.User.Password = ""

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: req.RefreshToken, // Reutilizamos el mismo
		User:         refreshToken.User,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout invalida la sesión del usuario usando estrategia híbrida
func (s *AuthService) Logout(token string, userID uint, role models.UserRole) error {
	// 1. SIEMPRE eliminar físicamente todos los refresh tokens del usuario
	err := s.refreshTokenRepo.DeleteByUserID(userID)
	if err != nil {
		return errors.New("error deleting refresh tokens")
	}

	// 2. SOLO para usuarios privilegiados (admin/teacher): añadir JWT a blacklist
	if utils.IsPrivilegedUser(role) {
		// Extraer claims del token para obtener JTI y expiración
		claims, err := s.jwtService.ValidateToken(token)
		if err != nil {
			// Si el token ya es inválido, no necesitamos blacklist
			// pero el refresh token ya fue eliminado, así que logout es exitoso
			return nil
		}

		// Añadir JWT a blacklist para invalidación inmediata
		blacklistEntry := &models.TokenBlacklist{
			JTI:       claims.ID,
			Token:     token,
			UserID:    userID,
			ExpiresAt: claims.ExpiresAt.Time,
		}

		err = s.blacklistRepo.Add(blacklistEntry)
		if err != nil {
			return errors.New("error adding token to blacklist")
		}
	}

	// Para usuarios normales (student): solo se eliminó el refresh token
	// El JWT expirará naturalmente en máximo 15 minutos sin posibilidad de renovación
	return nil
}

// IsTokenBlacklisted verifica si un token está en la blacklist
// Solo se debe llamar para usuarios privilegiados (admin/teacher)
func (s *AuthService) IsTokenBlacklisted(jti string, role models.UserRole) (bool, error) {
	// Solo verificar blacklist para usuarios privilegiados
	if !utils.IsPrivilegedUser(role) {
		return false, nil // Usuarios normales no usan blacklist
	}

	return s.blacklistRepo.IsBlacklisted(jti)
}