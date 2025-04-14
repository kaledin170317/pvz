package myhttp

import (
	"context"
	"net/http"
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
				WriteError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return am.secret, nil
			})
			if err != nil || !token.Valid {
				WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			role, ok := claims["role"].(string)
			userID, idOk := claims["id"].(string)
			email, emailOk := claims["email"].(string)
			if !ok || !idOk || !emailOk {
				WriteError(w, http.StatusUnauthorized, "invalid token payload")
				return
			}

			if role != requiredRole {
				WriteError(w, http.StatusForbidden, "insufficient role")
				return
			}

			ctx := context.WithValue(r.Context(), "id", userID)
			ctx = context.WithValue(r.Context(), "email", email)
			ctx = context.WithValue(ctx, "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
