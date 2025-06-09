package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"soccer-api/internal/api/apicontext"
	"soccer-api/internal/api/helpers"
	"soccer-api/internal/api/request"
	"soccer-api/internal/api/response"
	"soccer-api/internal/domain/entity"
	playerService "soccer-api/internal/domain/service/player"
	"soccer-api/internal/infrastructure/observer"
	localizationService "soccer-api/internal/localization"
)

const PLAYER_ID = "playerID"

type Player struct {
	playerService       *playerService.Player
	observer            *observer.Observer
	localizationService *localizationService.Service
}

func New(c Configuration) *Player {
	return &Player{
		playerService:       c.PlayerService,
		observer:            c.Observer,
		localizationService: c.LocalizationService,
	}
}

func (p *Player) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)

	credentials := apicontext.GetUser(ctx)
	if credentials == nil {
		msg := p.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.UnauthorizedResponse(w, msg)
		return
	}

	playerParam := chi.URLParam(r, PLAYER_ID)
	playerID, err := uuid.Parse(playerParam)
	if err != nil {
		msg := p.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	var req request.UpdatePlayerRequest

	if err := helpers.ReadJSON(w, r, &req); err != nil {
		msg := p.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Country = strings.TrimSpace(req.Country)

	updatedPlayer, err := p.playerService.UpdatePlayerDetails(ctx, credentials.ID, playerID, req.FirstName, req.LastName, req.Country)
	if err != nil {
		var msg string
		switch {
		case errors.Is(err, playerService.ErrUpdateNoFields):
			msg = p.localizationService.GetMessage(localizer, playerService.MsgPlayerUpdateNoFields)
			helpers.BadRequestResponse(w, msg)
		case errors.Is(err, playerService.ErrInvalidCountry):
			msg = p.localizationService.GetMessage(localizer, playerService.MsgInvalidCountry, map[string]string{"Value": req.Country})
			helpers.BadRequestResponse(w, msg)
		case errors.Is(err, playerService.ErrPlayerNotFound):
			msg = p.localizationService.GetMessage(localizer, playerService.MsgPlayerNotFound)
			helpers.NotFoundResponse(w, msg)
		case errors.Is(err, playerService.ErrPlayerNotOwned):
			msg = p.localizationService.GetMessage(localizer, playerService.MsgPlayerNotOwned)
			helpers.ErrorResponse(w, http.StatusForbidden, msg)
		default:
			msg = p.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
			helpers.ServerErrorResponse(w, msg)
		}
		return
	}

	respData := response.PlayerDetailResponse{
		ID:          updatedPlayer.ID,
		TeamID:      updatedPlayer.TeamID,
		FirstName:   updatedPlayer.FirstName,
		LastName:    updatedPlayer.LastName,
		Country:     entity.Country(p.localizationService.LocalizeCountry(localizer, entity.Country(updatedPlayer.Country))),
		Age:         updatedPlayer.Age,
		Position:    entity.PlayerPosition(p.localizationService.LocalizePlayerPosition(localizer, updatedPlayer.Position)),
		MarketValue: updatedPlayer.MarketValue,
		CreatedAt:   updatedPlayer.CreatedAt,
		UpdatedAt:   updatedPlayer.UpdatedAt,
	}

	helpers.WriteJSON(w, http.StatusOK, respData, nil)
}
