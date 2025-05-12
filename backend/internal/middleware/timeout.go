package middleware

import (
	"context"
	"net/http"
	"time"
)

// New timeout with context
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			// Pass the context with timeout to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
