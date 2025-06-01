package service

import (
	"soccer-api/internal/domain/repository"
)

type Configuration struct {
	AuthRepo  repository.Auth
	UserRepo  repository.User
	TokenRepo repository.Token
}
