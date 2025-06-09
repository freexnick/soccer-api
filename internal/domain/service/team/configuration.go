package service

import (
	"soccer-api/internal/domain/repository"
	countryService "soccer-api/internal/domain/service/country"
	playerService "soccer-api/internal/domain/service/player"
	txService "soccer-api/internal/domain/service/uow"
	userService "soccer-api/internal/domain/service/user"
)

type Configuration struct {
	TxService      txService.TxManager
	TeamRepo       repository.Team
	UserService    userService.User
	PlayerService  playerService.Player
	CountryService countryService.Country
}
