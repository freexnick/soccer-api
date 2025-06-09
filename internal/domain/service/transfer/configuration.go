package service

import (
	"soccer-api/internal/domain/repository"
	playerService "soccer-api/internal/domain/service/player"
	txService "soccer-api/internal/domain/service/uow"
)

type Configuration struct {
	TransferRepo  repository.Transfer
	TxService     txService.TxManager
	PlayerService playerService.Player
}
