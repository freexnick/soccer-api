package handlers

import (
	"errors"
	"net/http"

	"soccer-api/internal/api/helpers"
	"soccer-api/internal/api/request"
	"soccer-api/internal/api/response"
	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/service/auth"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/localization"
)

type Auth struct {
	AuthService         *service.Auth
	Observer            *observer.Observer
	LocalizationService *localization.Service
}

func New(c Configuration) *Auth {
	return &Auth{
		AuthService:         c.AuthService,
		Observer:            c.Observer,
		LocalizationService: c.LocalizationService,
	}
}

func (a *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.RegisterRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.BadRequestResponse(w, err)
		return
	}

	user, err := a.AuthService.Register(ctx, &entity.User{Email: req.Email, Password: req.Password}, "")
	if err != nil {
		helpers.BadRequestResponse(w, err)
		return
	}

	respUser := response.UserResponse{
		UserID: user.ID,
		Email:  user.Email,
		Token:  user.Token,
	}

	if err := helpers.WriteJSON(w, http.StatusCreated, respUser, nil); err != nil {
		a.Observer.Error(ctx, errors.New("AuthHandler.SignUp: WriteJSON for success response failed"))
	}
}

func (a *Auth) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.LoginRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.BadRequestResponse(w, err)
		return
	}

	user, err := a.AuthService.Login(ctx, &entity.User{Email: req.Email, Password: req.Password})
	if err != nil {
		helpers.BadRequestResponse(w, err)
		return
	}

	if err := helpers.WriteJSON(w, http.StatusOK, user, nil); err != nil {
		a.Observer.Error(ctx, errors.New("AuthHandler.SignIn: WriteJSON for success response failed"))
	}
}
