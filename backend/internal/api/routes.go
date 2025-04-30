package api

import (
	"mana/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Recover)
	router.Use(middleware.Logging)

	// get health endpoint
	router.Get("/api/v1/health", Health_Handler)

	return router
}
