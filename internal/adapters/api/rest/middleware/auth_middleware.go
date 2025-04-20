package middleware

import (
	"context"
	"net/http"
	"pvZ/internal/adapters/api/rest"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret []byte) *AuthMiddleware {
	return &AuthMiddleware{secret: secret}
}
func (am *AuthMiddleware) RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				rest.WriteError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return am.secret, nil
			})
			if err != nil || !token.Valid {
				rest.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			role, ok := claims["role"].(string)

			if !ok {
				rest.WriteError(w, http.StatusUnauthorized, "invalid token payload")
				return
			}

			if role != requiredRole {
				rest.WriteError(w, http.StatusForbidden, "insufficient role")
				return
			}

			ctx := context.WithValue(r.Context(), "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
