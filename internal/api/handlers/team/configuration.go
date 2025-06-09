package handlers

import (
	teamService "soccer-api/internal/domain/service/team"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/localization"
)

type Configuration struct {
	TeamService         *teamService.Team
	Observer            *observer.Observer
	LocalizationService *localization.Service
}
