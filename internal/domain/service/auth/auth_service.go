package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
	teamService "soccer-api/internal/domain/service/team"
	tokenService "soccer-api/internal/domain/service/token"
	txService "soccer-api/internal/domain/service/uow"
	userService "soccer-api/internal/domain/service/user"
)

type Auth struct {
	txService    txService.TxManager
	authRepo     repository.Auth
	userService  userService.User
	tokenService tokenService.Token
	teamService  teamService.Team
}

func New(c Configuration) *Auth {
	return &Auth{
		txService:    c.TxService,
		authRepo:     c.AuthRepo,
		userService:  c.UserService,
		tokenService: c.TokenService,
		teamService:  c.TeamService,
	}
}

func (a *Auth) Login(ctx context.Context, user *entity.User) (*entity.Credentials, error) {
	foundUser, err := a.userService.GetByEmail(ctx, strings.ToLower(user.Email))
	if err != nil {
		return nil, err
	}

	if foundUser == nil || !comparePasswords(user.Password, foundUser.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := a.tokenService.GenerateToken(ctx, *foundUser)
	if err != nil {
		return nil, ErrSessionCreationFailed
	}

	return &entity.Credentials{
		ID:    foundUser.ID,
		Email: foundUser.Email,
		Token: token,
	}, nil
}

func (a *Auth) Register(ctx context.Context, user *entity.User, token string) (*entity.Credentials, error) {
	if user.Email == "" || user.Password == "" {
		return nil, ErrValidationFailed
	}

	user.Email = strings.ToLower(user.Email)

	checkUser, err := a.userService.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if checkUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := hashPasswords(user.Password)
	if err != nil {
		return nil, err
	}
	user.ID = uuid.New()
	user.Password = hashedPassword

	uowErr := a.txService.Execute(ctx, func(repos repository.Repositories) error {
		if err := a.userService.Create(ctx, user); err != nil {
			return err
		}

		emailPrefix := strings.SplitN(user.Email, "@", 2)[0]
		emailPrefix = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(emailPrefix, ".", " "), "_", " "))
		emailPrefix = strings.ToUpper(emailPrefix[:1]) + emailPrefix[1:]

		team, err := a.teamService.CreateTeamAndInitialPlayers(ctx, repos, user.ID, emailPrefix)
		if err != nil {
			return err
		}

		user.Team = *team

		return nil
	})

	if uowErr != nil {
		return nil, err
	}

	generatedToken, err := a.tokenService.GenerateToken(ctx, *user)
	if err != nil {
		return nil, err
	}

	return &entity.Credentials{
		ID:    user.ID,
		Email: user.Email,
		Token: generatedToken,
	}, nil
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
