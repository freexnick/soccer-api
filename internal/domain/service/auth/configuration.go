package service

import (
	"soccer-api/internal/domain/repository"
	teamService "soccer-api/internal/domain/service/team"
	tokenService "soccer-api/internal/domain/service/token"
	txService "soccer-api/internal/domain/service/uow"
	userService "soccer-api/internal/domain/service/user"
)

type Configuration struct {
	TxService    txService.TxManager
	AuthRepo     repository.Auth
	UserService  userService.User
	TokenService tokenService.Token
	TeamService  teamService.Team
}
