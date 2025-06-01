package repository

import (
	"context"
	"soccer-api/internal/domain/entity"
)

type Auth interface {
	Register(ctx context.Context, user *entity.User, token string) (*entity.Credentials, error)
	Login(ctx context.Context, user *entity.User) (*entity.Credentials, error)
}
