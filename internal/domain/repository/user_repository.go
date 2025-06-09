package repository

import (
	"context"
	"soccer-api/internal/domain/entity"

	"github.com/google/uuid"
)

type User interface {
	BaseRepository[entity.User]
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}
