package middleware

import (
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware maneja la autenticación de requests
type AuthMiddleware struct {
	// TODO: Agregar dependencias como AuthService
}

// NewAuthMiddleware crea una nueva instancia del middleware de auth
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
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
		
		// TODO: Validar token con AuthService
		// user, err := m.authService.ValidateToken(token)
		// if err != nil {
		//     http.Error(w, "Invalid token", http.StatusUnauthorized)
		//     return
		// }

		// TODO: Agregar user al context
		// ctx := context.WithValue(r.Context(), "user", user)
		// r = r.WithContext(ctx)

		// Por ahora, solo verificar que el token no esté vacío
		if token == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

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