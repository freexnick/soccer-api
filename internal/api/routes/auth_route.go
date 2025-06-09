package routes

import (
	authHandler "soccer-api/internal/api/handlers/auth"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, handler *authHandler.Auth) chi.Router {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", handler.SignUp)
		r.Post("/signin", handler.SignIn)
	})
	return r
}
