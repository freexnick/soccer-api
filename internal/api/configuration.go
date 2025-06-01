package api

import (
	"github.com/go-chi/chi/v5"

	authHandler "soccer-api/internal/api/handlers/auth"
	authMW "soccer-api/internal/api/middlewares/auth"
	localizationMW "soccer-api/internal/api/middlewares/localization"
	observerMW "soccer-api/internal/api/middlewares/observer"
	authserviceMW "soccer-api/internal/domain/service/auth"
	tokenservice "soccer-api/internal/domain/service/token"
	userservice "soccer-api/internal/domain/service/user"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/localization"
)

type Configuration struct {
	APIVersion             string
	Router                 chi.Router
	Observer               *observer.Observer
	ObserverMiddleware     *observerMW.ObserverMiddleware
	AuthMiddleware         *authMW.AuthMiddleware
	LocalizationMiddleware *localizationMW.LocalizationMiddleware
	AuthHandler            *authHandler.Auth
	AuthService            *authserviceMW.Auth
	TokenService           *tokenservice.Token
	UserService            *userservice.User
	LocalizationService    *localization.Service
}
