package auth

import (
	"context"

	"gorm.io/gorm"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/infrastructure/database/gorm/models"
)

type Auth struct {
	modal models.UserModel
	db    *gorm.DB
}

func New(db *gorm.DB) *Auth {
	return &Auth{db: db}
}

func (a *Auth) Login(ctx context.Context, userCredentials *entity.User) (*entity.Credentials, error) {
	var currentUser models.UserModel

	if err := a.db.WithContext(ctx).Where("email = ?", userCredentials.Email).First(&currentUser).Error; err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *Auth) Register(ctx context.Context, user *entity.User, token string) (*entity.Credentials, error) {
	newUser := a.modal.ToUserModal(user)

	if err := a.db.WithContext(ctx).Create(&newUser).Error; err != nil {
		return nil, err
	}

	registeredUser := newUser.ToUserEntity()

	credentials := entity.Credentials{
		ID:    registeredUser.ID,
		Email: registeredUser.Email,
		Token: token,
	}

	return &credentials, nil
}
