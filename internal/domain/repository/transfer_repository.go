package repository

import (
	"context"
	"soccer-api/internal/domain/entity"

	"github.com/google/uuid"
)

type Transfer interface {
	BaseRepository[entity.TransferListing]
	GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*entity.TransferListing, error)
	GetAll(ctx context.Context) ([]entity.TransferListing, error)
}
