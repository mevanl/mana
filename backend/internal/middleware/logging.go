package middleware

import (
	"log"
	"net/http"
	"time"
)

// logs the HTTP method, path, response status, and duration for each request.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start_time := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)
		duration_time := time.Since(start_time)

		log.Printf("[%s] %s %s %d %s",
			start_time.Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			duration_time,
		)
	})
}

// Catches status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriterHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
