package service

import (
	"context"
	"errors"
	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"

	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	authRepo  repository.Auth
	userRepo  repository.User
	tokenRepo repository.Token
}

func New(c Configuration) *Auth {
	return &Auth{
		authRepo:  c.AuthRepo,
		userRepo:  c.UserRepo,
		tokenRepo: c.TokenRepo,
	}
}

func (a *Auth) Login(ctx context.Context, user *entity.User) (*entity.Credentials, error) {
	foundUser, err := a.userRepo.GetByEmail(ctx, strings.ToLower(user.Email))
	if err != nil {
		return nil, err
	}

	if foundUser == nil || !comparePasswords(user.Password, foundUser.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := a.tokenRepo.GenerateToken(ctx, *foundUser)
	if err != nil {
		return nil, err
	}

	return &entity.Credentials{
		ID:    foundUser.ID,
		Email: foundUser.Email,
		Token: token,
	}, nil
}

func (a *Auth) Register(ctx context.Context, user *entity.User, token string) (*entity.Credentials, error) {
	checkUser, err := a.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if checkUser != nil {
		return nil, errors.New("user already exist")
	}

	hashedPassword, err := hashPasswords(user.Password)
	if err != nil {
		return nil, err
	}

	user.Email = strings.ToLower(user.Email)
	user.Password = hashedPassword

	generatedToken, err := a.tokenRepo.GenerateToken(ctx, *user)
	if err != nil {
		return nil, err
	}

	return a.authRepo.Register(ctx, user, generatedToken)
}

func hashPasswords(plainPassword string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func comparePasswords(plainPassword, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}
