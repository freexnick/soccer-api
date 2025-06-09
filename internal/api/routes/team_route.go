package routes

import (
	"github.com/go-chi/chi/v5"

	teamHanler "soccer-api/internal/api/handlers/team"
)

func TeamRoutes(r chi.Router, handler *teamHanler.Team) chi.Router {
	r.Route("/team", func(r chi.Router) {
		r.Get("/", handler.GetMyTeam)
		r.Put("/", handler.UpdateMyTeam)
	})
	return r
}
