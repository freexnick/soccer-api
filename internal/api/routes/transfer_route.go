package routes

import (
	"github.com/go-chi/chi/v5"

	transferHandler "soccer-api/internal/api/handlers/transfer"
)

func TransferRoutes(r chi.Router, handler *transferHandler.Transfer) chi.Router {
	r.Route("/transfers", func(r chi.Router) {
		r.Get("/", handler.ViewTransferList)
		r.Group(func(r chi.Router) {
			r.Route("/{playerID}", func(r chi.Router) {
				r.Post("/", handler.ListPlayerForTransfer)
				r.Delete("/", handler.RemovePlayerFromTransfer)
			})
		})

		r.Post("/buy/{playerID}", handler.BuyPlayer)
	})
	return r
}
