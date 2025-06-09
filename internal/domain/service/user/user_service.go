package service

import (
	"context"

	"github.com/google/uuid"

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

func (u *User) GetByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.repo.GetByID(ctx, userID)
}

func (u *User) Create(ctx context.Context, user *entity.User) error {
	return u.repo.Create(ctx, user)
}
