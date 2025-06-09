package routes

import (
	"github.com/go-chi/chi/v5"
	playerHandler "soccer-api/internal/api/handlers/player"
)

func PlayerRoutes(r chi.Router, handler *playerHandler.Player) chi.Router {
	r.Route("/player", func(r chi.Router) {
		r.Put("/{playerID}", handler.UpdatePlayer)
	})
	return r
}
