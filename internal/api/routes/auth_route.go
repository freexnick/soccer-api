package routes

import (
	authHandler "soccer-api/internal/api/handlers/auth"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, authHandler *authHandler.Auth) chi.Router {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.SignUp)
		r.Post("/signin", authHandler.SignIn)
	})

	return r
}
