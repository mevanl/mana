package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	// middleware pipeline:
	// recover > log > security header > CORS > IP Filter > Authenticate >
	// Authorize > Rate Limit > Input Validation > Content Moderation > AntiSpam/Abuse >
	// Session Management > Encrypt/Integrity > Audit Log > Timeout

	// add our middlewares
	// router.Use(middlewareRecover)
	router.Use(middlewareLogger)

	// get health endpoint
	router.Get("/api/v1/health", Health_Handler)

	return router
}

// middlewareLogger logs a request with: [timestamp_start] [http method] [request_url] [duration]
func middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start_time := time.Now()
		next.ServeHTTP(w, r)
		duration_time := time.Since(start_time)

		log.Printf("[%s] %s %s %s", start_time.Format(time.RFC3339), r.Method, r.URL.Path, duration_time)
	})
}

// func middlewareRecover(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// catch panics and handle errors
// 		defer func() {
// 			if recover := recover(); recover != nil {
// 				// log panic error
// 				log.Printf("Panic recovered: %+v", recover)

// 				// return 500 (internal server error)
// 				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 			}
// 		}()

// 		// call next handler in chain
// 		next.ServeHTTP(w, r)
// 	})
// }
