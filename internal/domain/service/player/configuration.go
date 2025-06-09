package service

import (
	"soccer-api/internal/domain/repository"
	countryService "soccer-api/internal/domain/service/country"
	randomService "soccer-api/internal/domain/service/random"
	txService "soccer-api/internal/domain/service/uow"
)

type Configuration struct {
	PlayerRepo     repository.Player
	TxService      txService.TxManager
	RandomService  randomService.Random
	CountryService countryService.Country
}
