package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/render"
	"llm-promp-inj.api/internal/service"
)

type AuthMiddleware struct {
	tokenService *service.TokenService
	authService  *service.AuthenticationService
}

func NewAuthMiddleware(tokenService *service.TokenService, authService *service.AuthenticationService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService, authService: authService}
}

func (m *AuthMiddleware) Authenticate(requiredRole []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"status": "Unauthorized"})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"status": "Token Error"})
				return
			}

			tokenString := parts[1]

			claims, err := m.tokenService.ExtractAndValidateToken(tokenString)

			if claims == nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"status": "Unauthorized"})
				return
			}

			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{"status": "Internal Server Error"})
				return
			}

			if !m.authService.IsValidSession(claims.SessionID, claims.Sub) {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"status": "Unauthorized"})
				return
			}

			// Enforce role-based authorization
			if len(requiredRole) > 0 && !slices.Contains(requiredRole, claims.Data["role"]) {
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, map[string]string{"status": "Forbidden"})
				return
			}

			ctx := context.WithValue(r.Context(), "userClaims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
