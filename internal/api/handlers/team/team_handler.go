package handlers

import (
	"errors"
	"net/http"
	"strings"

	"soccer-api/internal/api/apicontext"
	"soccer-api/internal/api/helpers"
	"soccer-api/internal/api/request"
	"soccer-api/internal/api/response"
	"soccer-api/internal/domain/entity"
	teamService "soccer-api/internal/domain/service/team"
	"soccer-api/internal/infrastructure/observer"
	localizationService "soccer-api/internal/localization"
)

type Team struct {
	teamService         *teamService.Team
	observer            *observer.Observer
	localizationService *localizationService.Service
}

func New(c Configuration) *Team {
	return &Team{
		teamService:         c.TeamService,
		observer:            c.Observer,
		localizationService: c.LocalizationService,
	}
}

func (t *Team) GetMyTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)
	credentials := apicontext.GetUser(ctx)

	if credentials == nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.UnauthorizedResponse(w, msg)
		return
	}

	team, totalValue, err := t.teamService.GetTeamByUser(ctx, credentials.ID)
	if err != nil {
		var msg string
		if errors.Is(err, teamService.ErrUserHasNoTeam) {
			msg = t.localizationService.GetMessage(localizer, teamService.MsgTeamNotFound)
			helpers.NotFoundResponse(w, msg)
			return
		}
		msg = t.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
		helpers.NotFoundResponse(w, msg)
		return
	}

	playerDTOs := make([]response.PlayerBasicResponse, len(team.Players))

	for i, p := range team.Players {
		playerDTOs[i] = response.PlayerBasicResponse{
			ID:          p.ID,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			Position:    entity.PlayerPosition(t.localizationService.LocalizePlayerPosition(localizer, p.Position)),
			Age:         p.Age,
			Country:     entity.Country(t.localizationService.LocalizeCountry(localizer, entity.Country(p.Country))),
			MarketValue: p.MarketValue,
		}
	}

	resp := response.TeamResponse{
		ID:         team.ID,
		UserID:     team.UserID,
		Name:       team.Name,
		Country:    entity.Country(t.localizationService.LocalizeCountry(localizer, team.Country)),
		Budget:     team.Budget,
		TotalValue: totalValue,
		Players:    playerDTOs,
		CreatedAt:  team.CreatedAt,
		UpdatedAt:  team.UpdatedAt,
	}

	helpers.WriteJSON(w, http.StatusOK, resp, nil)
}

func (t *Team) UpdateMyTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)

	credentials := apicontext.GetUser(ctx)
	if credentials == nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.UnauthorizedResponse(w, msg)
		return
	}

	var req request.UpdateTeamRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Country = strings.TrimSpace(req.Country)

	team, totalValue, err := t.teamService.Update(ctx, credentials.ID, req.Name, req.Country)
	if err != nil {
		var msg string
		switch {
		case errors.Is(err, teamService.ErrUpdateNoFields):
			msg = t.localizationService.GetMessage(localizer, teamService.MsgTeamUpdateNoFields)
			helpers.BadRequestResponse(w, msg)
		case errors.Is(err, teamService.ErrInvalidCountry):
			msg = t.localizationService.GetMessage(localizer, teamService.MsgInvalidCountry, map[string]string{"Value": req.Country})
			helpers.BadRequestResponse(w, msg)
		case errors.Is(err, teamService.ErrUserHasNoTeam):
			msg = t.localizationService.GetMessage(localizer, teamService.MsgUserHasNoTeam)
			helpers.NotFoundResponse(w, msg)
		default:
			t.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
			helpers.ServerErrorResponse(w, msg)
		}
		return
	}

	playerDTOs := make([]response.PlayerBasicResponse, len(team.Players))
	for i, p := range team.Players {
		playerDTOs[i] = response.PlayerBasicResponse{
			ID:          p.ID,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			Position:    entity.PlayerPosition(t.localizationService.LocalizePlayerPosition(localizer, p.Position)),
			Country:     entity.Country(t.localizationService.LocalizeCountry(localizer, entity.Country(p.Country))),
			Age:         p.Age,
			MarketValue: p.MarketValue,
		}
	}

	respData := response.TeamResponse{
		ID:         team.ID,
		UserID:     team.UserID,
		Name:       team.Name,
		Country:    entity.Country(t.localizationService.LocalizeCountry(localizer, team.Country)),
		Budget:     team.Budget,
		TotalValue: totalValue,
		Players:    playerDTOs,
		CreatedAt:  team.CreatedAt,
		UpdatedAt:  team.UpdatedAt,
	}

	helpers.WriteJSON(w, http.StatusOK, respData, nil)

}
