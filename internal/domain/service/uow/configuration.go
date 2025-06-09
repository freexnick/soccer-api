package service

import "soccer-api/internal/domain/repository"

type Configuration struct {
	TxManager repository.TransactionManager
}
