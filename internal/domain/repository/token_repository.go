package repository

import (
	"context"
	"net/http"
	"soccer-api/internal/domain/entity"
)

type Token interface {
	GenerateToken(ctx context.Context, user entity.User) (string, error)
	ValidateToken(ctx context.Context, req *http.Request) (*entity.Token, error)
}
