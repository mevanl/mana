package middleware

import (
	"log"
	"net/http"
	"time"
)

// logs a request with: [timestamp_start] [http method] [request_url] [duration]
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start_time := time.Now()
		next.ServeHTTP(w, r)
		duration_time := time.Since(start_time)

		log.Printf("[%s] %s %s %s", start_time.Format(time.RFC3339), r.Method, r.URL.Path, duration_time)
	})
}
