package repository

import (
	"context"

	"github.com/google/uuid"

	"soccer-api/internal/domain/entity"
)

type Team interface {
	BaseRepository[entity.Team]
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*entity.Team, error)
	GetTeamByUser(ctx context.Context, userID uuid.UUID) (*entity.Team, error)
}
