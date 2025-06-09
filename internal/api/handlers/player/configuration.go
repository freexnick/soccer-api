package handlers

import (
	playerService "soccer-api/internal/domain/service/player"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/localization"
)

type Configuration struct {
	PlayerService       *playerService.Player
	Observer            *observer.Observer
	LocalizationService *localization.Service
}
