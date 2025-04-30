package middleware

import (
	"net/http"
	"time"
)

func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			http.TimeoutHandler(next, duration, `{"error": "request timed out"}`).ServeHTTP(w, r)
		})
	}
}
