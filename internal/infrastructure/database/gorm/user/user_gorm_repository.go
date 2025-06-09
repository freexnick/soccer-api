package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/infrastructure/database/gorm/models"
)

type User struct {
	model models.UserModel
	db    *gorm.DB
}

func New(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, user *entity.User) error {
	userModel := u.model.ToUserModal(user)

	if err := u.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel models.UserModel

	if err := u.db.WithContext(ctx).
		Preload("Team.Players").
		Where("email = ?", email).
		First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return userModel.ToUserEntity(), nil
}

func (u *User) GetByID(ctx context.Context, ID uuid.UUID) (*entity.User, error) {
	var userModel models.UserModel

	if err := u.db.WithContext(ctx).
		Preload("Team.Players").
		Where("id = ?", ID).
		First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return userModel.ToUserEntity(), nil
}

func (u *User) Update(ctx context.Context, user *entity.User) error {
	return nil
}

func (u *User) Delete(ctx context.Context, userID uuid.UUID) error {
	return nil
}
