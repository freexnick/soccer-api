package repository

import (
	"context"
	"soccer-api/internal/domain/entity"
)

type User interface {
	BaseRepository[entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
