package handlers

import (
	service "soccer-api/internal/domain/service/auth"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/localization"
)

type Configuration struct {
	AuthService         *service.Auth
	Observer            *observer.Observer
	LocalizationService *localization.Service
}
