package middleware

import (
	"context"
	"net/http"
	"strings"

	"mana/internal/auth"
)

// key type avoids collisions in context
type contextKey string

const UserIDKey contextKey = "userID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{ "error": "missing or invalid authorization header" }`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, `{ "error": "invalid or expired token" }`, http.StatusUnauthorized)
			return
		}

		// Store userID in context for downstream access
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
