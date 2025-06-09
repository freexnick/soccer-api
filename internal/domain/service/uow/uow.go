package service

import (
	"context"

	"soccer-api/internal/domain/repository"
)

type TxManager struct {
	repo repository.TransactionManager
}

func New(c Configuration) *TxManager {
	return &TxManager{repo: c.TxManager}
}

func (tm *TxManager) Execute(ctx context.Context, fn repository.RepoFunc) error {
	return tm.repo.Execute(ctx, fn)
}
