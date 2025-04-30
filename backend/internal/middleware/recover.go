package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// catch panics and handle errors
		defer func() {
			if recover := recover(); recover != nil {
				// log panic error
				log.Printf("Panic recovered: %+v\nRequest: %s %s\nStack trace:\n%s",
					recover, r.Method, r.URL.Path, string(debug.Stack()))

				// return 500 (internal server error)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		// call next handler in chain
		next.ServeHTTP(w, r)
	})
}
