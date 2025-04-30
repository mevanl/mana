package middleware

import (
	"log"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// catch panics and handle errors
		defer func() {
			if recover := recover(); recover != nil {
				// log panic error
				log.Printf("Panic recovered: %+v", recover)

				// return 500 (internal server error)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		// call next handler in chain
		next.ServeHTTP(w, r)
	})
}
