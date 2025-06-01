package user

import (
	"context"
	"errors"
	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type user struct {
	model UserModel
	db    *gorm.DB
}

func New(db *gorm.DB) repository.User {
	return &user{db: db}
}

func (u *user) Create(ctx context.Context, user *entity.User) error {
	gormModel := u.model.FromDomain(user)

	if err := u.db.WithContext(ctx).Create(&gormModel).Error; err != nil {
		return err
	}

	user.CreatedAt = gormModel.CreatedAt
	user.UpdatedAt = gormModel.UpdatedAt

	return nil
}

func (u *user) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var gormModel UserModel

	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&gormModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return gormModel.ToDomain(), nil
}

func (u *user) Update(ctx context.Context, dUser *entity.User) error {
	return errors.New("not implemented")
}
func (u *user) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}
func (u *user) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return nil, errors.New("not implemented")
}
