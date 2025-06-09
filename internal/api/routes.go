package api

import (
	"net/http"
	"soccer-api/internal/api/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func New(conf Configuration) http.Handler {
	r := conf.Router
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSFR-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(conf.ObserverMiddleware.Observe)
	r.Use(conf.LocalizationMiddleware.Localize)

	r.Route(conf.APIVersion, func(r chi.Router) {
		routes.AuthRoutes(r, conf.AuthHandler)

		r.Group(func(r chi.Router) {
			r.Use(conf.AuthMiddleware.Authenticate)
			routes.TeamRoutes(r, conf.TeamHandler)
			routes.PlayerRoutes(r, conf.PlayerHandler)
			routes.TransferRoutes(r, conf.TransferHandler)
		})
	})

	return r
}
