package service

import (
	"context"
	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
)

type User struct {
	repo repository.User
}

func New(repo repository.User) *User {
	return &User{repo: repo}
}

func (u *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.repo.GetByEmail(ctx, email)
}
