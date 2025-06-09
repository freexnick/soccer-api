package repository

import (
	"context"

	"soccer-api/internal/domain/entity"

	"github.com/google/uuid"
)

type Player interface {
	BaseRepository[entity.Player]
	CreateBatch(ctx context.Context, players []entity.Player, teamID uuid.UUID) error
	GetByID(ctx context.Context, playerID uuid.UUID) (*entity.Player, error)
}
