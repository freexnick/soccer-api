package service

import (
	"context"
	"net/http"
	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
)

type Token struct {
	repo repository.Token
}

func New(repo repository.Token) *Token {
	return &Token{repo: repo}
}

func (t *Token) GenerateToken(ctx context.Context, user entity.User) (string, error) {
	return t.repo.GenerateToken(ctx, user)
}

func (t *Token) ValidateToken(ctx context.Context, r *http.Request) (*entity.Token, error) {
	return t.repo.ValidateToken(ctx, r)
}
