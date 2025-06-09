package handlers

import (
	"errors"
	"net/http"

	"soccer-api/internal/api/apicontext"
	"soccer-api/internal/api/helpers"
	"soccer-api/internal/api/request"
	"soccer-api/internal/api/response"
	"soccer-api/internal/domain/entity"
	authService "soccer-api/internal/domain/service/auth"
	"soccer-api/internal/infrastructure/observer"
	localizationService "soccer-api/internal/localization"
)

type Auth struct {
	authService         *authService.Auth
	observer            *observer.Observer
	localizationService *localizationService.Service
}

func New(c Configuration) *Auth {
	return &Auth{
		authService:         c.AuthService,
		observer:            c.Observer,
		localizationService: c.LocalizationService,
	}
}

func (a *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)

	var req request.RegisterRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		msg := a.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	user, err := a.authService.Register(ctx, &entity.User{Email: req.Email, Password: req.Password}, "")
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrValidationFailed):
			msg := a.localizationService.GetMessage(localizer, authService.MsgValidationFailed)
			helpers.BadRequestResponse(w, msg)

		case errors.Is(err, authService.ErrUserAlreadyExists):
			msg := a.localizationService.GetMessage(localizer, authService.MsgUserAlreadyExists)
			helpers.BadRequestResponse(w, msg)
		default:
			msg := a.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
			helpers.ServerErrorResponse(w, msg)
		}
		return
	}

	respUser := response.UserResponse{
		UserID: user.ID,
		Email:  user.Email,
		Token:  user.Token,
	}

	helpers.WriteJSON(w, http.StatusCreated, respUser, nil)
}

func (a *Auth) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	localizer := apicontext.GetLocalizer(ctx)

	var req request.LoginRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		msg := a.localizationService.GetMessage(localizer, localizationService.MsgBadRequestError)
		helpers.BadRequestResponse(w, msg)
		return
	}

	user, err := a.authService.Login(ctx, &entity.User{Email: req.Email, Password: req.Password})
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidCredentials):
			msg := a.localizationService.GetMessage(localizer, authService.MsgInvalidCredentials)
			helpers.BadRequestResponse(w, msg)
		default:
			msg := a.localizationService.GetMessage(localizer, localizationService.MsgInternalServerError)
			helpers.ServerErrorResponse(w, msg)
		}
		return
	}

	helpers.WriteJSON(w, http.StatusOK, user, nil)
}
