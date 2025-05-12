package api

import (
	"mana/internal/middleware"
	"mana/internal/store"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewRouter(store *store.Store) http.Handler {
	router := chi.NewRouter()
	api := &API{Store: store}

	// Middleware
	router.Use(middleware.Recover)
	router.Use(middleware.Logging)
	router.Use(middleware.CORS)
	router.Use(middleware.SecurityHeaders)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.Authenticate)

	// Public routes
	router.Get("/api/v1/health", api.Health)
	router.Post("/api/v1/register", api.Register)
	router.Post("/api/v1/login", api.Login)

	// authenticated routes
	router.Route("/api/v1", func(r chi.Router) {
		// Guild
		r.Get("/guild/{id}", api.GetGuildByID)
		r.Get("/guilds", api.GetUserGuilds)
		r.Post("/guilds", api.CreateGuild)
		r.Delete("guild/{id}", api.DeleteGuild)
		r.Post("/guilds/invites/{code}", api.JoinGuildByInvite)

		// Channel
		r.Get("/guilds/{id}/channels", api.GetGuildChannels)

		// Messages
		r.Get("/channel/{id}/messages", api.GetMessagesByChannel)
		r.Post("/channel/{id}/messages", api.CreateMessage)
	})

	return router
}
