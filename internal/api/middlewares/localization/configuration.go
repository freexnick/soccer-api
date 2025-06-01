package middleware

import "soccer-api/internal/localization"

type Configuration struct {
	LocalizationService *localization.Service
}
