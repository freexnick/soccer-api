package handlers

import (
	transferService "soccer-api/internal/domain/service/transfer"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/localization"
)

type Configuration struct {
	TransferService     *transferService.Transfer
	Observer            *observer.Observer
	LocalizationService *localization.Service
}
