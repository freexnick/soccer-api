package middleware

import (
	tokenService "soccer-api/internal/domain/service/token"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/localization"
)

type Configuration struct {
	TokenService        tokenService.Token
	Observer            *observer.Observer
	LocalizationService *localization.Service
}
