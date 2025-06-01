package middleware

import (
	tokenService "soccer-api/internal/domain/service/token"
	"soccer-api/internal/infrastructure/observer"
)

type Configuration struct {
	TokenService tokenService.Token
	Observer     *observer.Observer
}
