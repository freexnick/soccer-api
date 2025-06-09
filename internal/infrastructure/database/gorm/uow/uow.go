package uow

import (
	"context"

	"gorm.io/gorm"

	"soccer-api/internal/domain/repository"
	"soccer-api/internal/infrastructure/database/gorm/player"
	"soccer-api/internal/infrastructure/database/gorm/team"
	"soccer-api/internal/infrastructure/database/gorm/transfer"
	"soccer-api/internal/infrastructure/database/gorm/user"
)

type TransactionManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) Execute(ctx context.Context, fn repository.RepoFunc) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transactionalRepos := repository.Repositories{
			User:     user.New(tx),
			Team:     team.New(tx),
			Player:   player.New(tx),
			Transfer: transfer.New(tx),
		}

		return fn(transactionalRepos)
	})
}
