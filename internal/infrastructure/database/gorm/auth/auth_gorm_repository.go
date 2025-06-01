package auth

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
	"soccer-api/internal/infrastructure/database/gorm/user"
)

type auth struct {
	modal user.UserModel
	db    *gorm.DB
}

func New(db *gorm.DB) repository.Auth {
	return &auth{db: db}
}

func (a *auth) Login(ctx context.Context, userCredentials *entity.User) (*entity.Credentials, error) {
	var currentUser user.UserModel

	if err := a.db.WithContext(ctx).Where("email = ?", userCredentials.Email).First(&currentUser).Error; err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *auth) Register(ctx context.Context, user *entity.User, token string) (*entity.Credentials, error) {
	user.ID = uuid.New()
	newUser := a.modal.FromDomain(user)

	if err := a.db.WithContext(ctx).Create(&newUser).Error; err != nil {
		return nil, err
	}

	registeredUser := newUser.ToDomain()

	credentials := entity.Credentials{
		ID:    registeredUser.ID,
		Email: registeredUser.Email,
		Token: token,
	}

	return &credentials, nil
}
