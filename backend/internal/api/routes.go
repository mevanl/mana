package api

import (
	"mana/internal/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Recover)
	router.Use(middleware.Logging)
	router.Use(middleware.CORS)
	router.Use(middleware.SecurityHeaders)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.Authenticate)

	// get health endpoint
	router.Get("/api/v1/health", Health_Handler)

	return router
}
