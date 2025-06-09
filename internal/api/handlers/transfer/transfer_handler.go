package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"soccer-api/internal/api/apicontext"
	"soccer-api/internal/api/helpers"
	"soccer-api/internal/api/request"
	"soccer-api/internal/api/response"
	"soccer-api/internal/domain/entity"
	playerService "soccer-api/internal/domain/service/player"
	transferService "soccer-api/internal/domain/service/transfer"
	"soccer-api/internal/infrastructure/observer"
	localizationService "soccer-api/internal/localization"
)

const PLAYER_ID = "playerID"

type Transfer struct {
	transferService     *transferService.Transfer
	observer            *observer.Observer
	localizationService *localizationService.Service
}

func New(c Configuration) *Transfer {
	return &Transfer{
		transferService:     c.TransferService,
		observer:            c.Observer,
		localizationService: c.LocalizationService,
	}
}

func (t *Transfer) ListPlayerForTransfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)

	credentials := apicontext.GetUser(ctx)
	if credentials == nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.UnauthorizedResponse(w, msg)
		return
	}

	playerIDStr := chi.URLParam(r, PLAYER_ID)
	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	var req request.ListPlayerForTransferRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	createdListing, err := t.transferService.ListPlayer(ctx, credentials.ID, playerID, req.AskingPrice)
	if err != nil {
		var msg string
		switch {
		case errors.Is(err, transferService.ErrPriceRequired):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgPriceRequired)
			helpers.BadRequestResponse(w, msg)
		case errors.Is(err, playerService.ErrPlayerNotFound):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgPlayerNotListed)
			helpers.NotFoundResponse(w, msg)
		case errors.Is(err, transferService.ErrListingPermissionDenied):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgListingPermissionDenied)
			helpers.ErrorResponse(w, http.StatusForbidden, msg)
		case errors.Is(err, transferService.ErrPlayerAlreadyListed):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgPlayerAlreadyListed)
			helpers.ErrorResponse(w, http.StatusConflict, msg)
		default:
			msg = t.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
			helpers.ServerErrorResponse(w, msg)
		}
		return
	}

	respData := response.PlayerListedResponse{
		Message:     t.localizationService.GetMessage(localizer, transferService.MsgListSuccessful),
		ListingID:   createdListing.ID,
		PlayerID:    createdListing.PlayerID,
		AskingPrice: createdListing.AskingPrice,
	}

	helpers.WriteJSON(w, http.StatusCreated, respData, nil)
}

func (t *Transfer) RemovePlayerFromTransfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	localizer := apicontext.GetLocalizer(ctx)

	credentials := apicontext.GetUser(ctx)
	if credentials == nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.UnauthorizedResponse(w, msg)
		return
	}

	playerIDStr := chi.URLParam(r, PLAYER_ID)
	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	serviceErr := t.transferService.RemovePlayerFromTransfer(ctx, credentials.ID, playerID)
	if serviceErr != nil {
		var msg string
		switch {
		case errors.Is(err, transferService.ErrPlayerNotListed):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgPlayerNotListed)
			helpers.NotFoundResponse(w, msg)
		case errors.Is(err, transferService.ErrListingPermissionDenied):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgListingPermissionDenied)
			helpers.ErrorResponse(w, http.StatusForbidden, msg)
		default:
			msg = t.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
			helpers.ServerErrorResponse(w, msg)
		}
		return
	}

	helpers.WriteJSON(w, http.StatusOK, nil, nil)
}

func (t *Transfer) ViewTransferList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)

	credentials := apicontext.GetUser(ctx)
	if credentials == nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.UnauthorizedResponse(w, msg)
		return
	}

	listings, err := t.transferService.ViewTransferList(ctx)
	if err != nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
		helpers.ServerErrorResponse(w, msg)
		return
	}

	resp := make([]response.TransferListItemResponse, len(listings))
	for i, listing := range listings {

		playerSummary := response.PlayerSummaryForListing{
			ID:          listing.Player.ID,
			FirstName:   listing.Player.FirstName,
			LastName:    listing.Player.LastName,
			Country:     entity.Country(t.localizationService.LocalizeCountry(localizer, entity.Country(listing.Player.Country))),
			Age:         listing.Player.Age,
			Position:    entity.PlayerPosition(t.localizationService.LocalizePlayerPosition(localizer, listing.Player.Position)),
			MarketValue: listing.Player.MarketValue,
		}

		teamSummary := response.TeamSummaryForListing{
			ID:      listing.SellingTeam.ID,
			Name:    listing.SellingTeam.Name,
			Country: entity.Country(t.localizationService.LocalizeCountry(localizer, entity.Country(listing.SellingTeam.Country))),
		}

		resp[i] = response.TransferListItemResponse{
			ListingID:   listing.ID,
			Player:      playerSummary,
			SellingTeam: teamSummary,
			AskingPrice: listing.AskingPrice,
			ListedAt:    listing.ListedAt,
		}
	}

	helpers.WriteJSON(w, http.StatusOK, resp, nil)
}

func (t *Transfer) BuyPlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)

	credentials := apicontext.GetUser(ctx)
	if credentials == nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgUnauthorizedError)
		helpers.UnauthorizedResponse(w, msg)
		return
	}

	playerIDStr := chi.URLParam(r, PLAYER_ID)
	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		msg := t.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	boughtPlayer, err := t.transferService.BuyPlayer(ctx, credentials.ID, playerID)
	if err != nil {
		var msg string
		switch {
		case errors.Is(err, transferService.ErrPlayerNotListed):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgPlayerNotListed)
			helpers.NotFoundResponse(w, msg)
		case errors.Is(err, playerService.ErrPlayerNotFound):
			msg = t.localizationService.GetMessage(localizer, playerService.MsgPlayerNotFound)
			helpers.NotFoundResponse(w, msg)
		case errors.Is(err, transferService.ErrCannotBuyOwn):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgCannotBuyOwn)
			helpers.ErrorResponse(w, http.StatusForbidden, msg)
		case errors.Is(err, transferService.ErrNotEnoughBudget):
			msg = t.localizationService.GetMessage(localizer, transferService.MsgNotEnoughBudget)
			helpers.ErrorResponse(w, http.StatusPaymentRequired, msg)
		default:
			msg = t.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
			helpers.ServerErrorResponse(w, msg)
		}
		return
	}

	respPlayer := response.PlayerDetailResponse{
		ID:          boughtPlayer.ID,
		TeamID:      boughtPlayer.TeamID,
		FirstName:   boughtPlayer.FirstName,
		LastName:    boughtPlayer.LastName,
		Country:     entity.Country(t.localizationService.LocalizeCountry(localizer, entity.Country(boughtPlayer.Country))),
		Age:         boughtPlayer.Age,
		Position:    entity.PlayerPosition(t.localizationService.LocalizePlayerPosition(localizer, boughtPlayer.Position)),
		MarketValue: boughtPlayer.MarketValue,
		CreatedAt:   boughtPlayer.CreatedAt,
		UpdatedAt:   boughtPlayer.UpdatedAt,
	}
	respData := response.PlayerPurchasedResponse{
		Message: t.localizationService.GetMessage(localizer, transferService.MsgPurchaseSuccessful),
		Player:  respPlayer,
	}

	helpers.WriteJSON(w, http.StatusOK, respData, nil)
}
