package api

import (
	"mana/internal/db"
	"mana/internal/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewRouter(store *db.Store) http.Handler {
	router := chi.NewRouter()
	api := &API{Store: store}

	// Middleware
	router.Use(middleware.Recover)
	router.Use(middleware.Logging)
	router.Use(middleware.CORS)
	router.Use(middleware.SecurityHeaders)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.Authenticate)

	router.Get("/api/v1/health", api.Health)
	router.Post("/api/v1/register", api.Register)
	router.Post("/api/v1/login", api.Login)

	return router
}
