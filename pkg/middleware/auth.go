package middleware

import (
	"context"
	"net/http"
	"strings"
	"backend-academi/internal/auth"
	"backend-academi/pkg/utils"
)

// AuthMiddleware maneja la autenticación de requests
type AuthMiddleware struct {
	authService *auth.AuthService
	jwtService  *utils.JWTService
}

// NewAuthMiddleware crea una nueva instancia del middleware de auth
func NewAuthMiddleware(authService *auth.AuthService, jwtService *utils.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		jwtService:  jwtService,
	}
}

// RequireAuth middleware que requiere autenticación
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Verificar formato Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		
		// Validar token JWT y extraer claims
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Verificar blacklist para usuarios privilegiados
		if utils.IsPrivilegedUser(claims.Role) {
			isBlacklisted, err := m.authService.IsTokenBlacklisted(claims.ID, claims.Role)
			if err != nil {
				http.Error(w, "Token validation error", http.StatusInternalServerError)
				return
			}
			if isBlacklisted {
				http.Error(w, "Token has been revoked", http.StatusUnauthorized)
				return
			}
		}

		// Agregar claims al context para usar en handlers
		ctx := context.WithValue(r.Context(), "user_claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// RequireRole middleware que requiere un role específico
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Obtener user del context y verificar role
			// user := r.Context().Value("user")
			// if user == nil || user.Role != role {
			//     http.Error(w, "Insufficient permissions", http.StatusForbidden)
			//     return
			// }

			next.ServeHTTP(w, r)
		})
	}
}