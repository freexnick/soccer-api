package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"soccer-api/internal/api"
	authHandler "soccer-api/internal/api/handlers/auth"
	playerHandler "soccer-api/internal/api/handlers/player"
	teamHandler "soccer-api/internal/api/handlers/team"
	transferHandler "soccer-api/internal/api/handlers/transfer"
	authMW "soccer-api/internal/api/middlewares/auth"
	localMW "soccer-api/internal/api/middlewares/localization"
	observMW "soccer-api/internal/api/middlewares/observer"
	token "soccer-api/internal/auth"
	"soccer-api/internal/config"
	"soccer-api/internal/domain/repository"
	authService "soccer-api/internal/domain/service/auth"
	coutnryService "soccer-api/internal/domain/service/country"
	playerService "soccer-api/internal/domain/service/player"
	randomService "soccer-api/internal/domain/service/random"
	teamService "soccer-api/internal/domain/service/team"
	tokenService "soccer-api/internal/domain/service/token"
	transferService "soccer-api/internal/domain/service/transfer"
	txService "soccer-api/internal/domain/service/uow"
	userService "soccer-api/internal/domain/service/user"
	"soccer-api/internal/infrastructure/database"
	"soccer-api/internal/infrastructure/database/gorm/auth"
	"soccer-api/internal/infrastructure/database/gorm/player"
	"soccer-api/internal/infrastructure/database/gorm/team"
	"soccer-api/internal/infrastructure/database/gorm/transfer"
	"soccer-api/internal/infrastructure/database/gorm/uow"
	"soccer-api/internal/infrastructure/database/gorm/user"
	"soccer-api/internal/infrastructure/observer"
	"soccer-api/internal/infrastructure/observer/logger"
	"soccer-api/internal/infrastructure/server"
	"soccer-api/internal/localization"
)

type Application struct {
	observ  *observer.Observer
	handler http.Handler
	server  server.Lifecycle
	db      database.Client

	txRepo             repository.TransactionManager
	userRepository     repository.User
	authRepository     repository.Auth
	tokenRepository    repository.Token
	teamRepository     repository.Team
	playerRepository   repository.Player
	transferRepository repository.Transfer

	txService           *txService.TxManager
	userService         *userService.User
	authService         *authService.Auth
	tokenService        *tokenService.Token
	teamService         *teamService.Team
	playerService       *playerService.Player
	countryService      *coutnryService.Country
	randomService       *randomService.Random
	localizationService *localization.Service
	transferService     *transferService.Transfer
	authMW              *authMW.AuthMiddleware
	observMW            *observMW.ObserverMiddleware
	localMW             *localMW.LocalizationMiddleware

	authHandler     *authHandler.Auth
	teamHandler     *teamHandler.Team
	playerHandler   *playerHandler.Player
	transferHandler *transferHandler.Transfer
}

func New(ctx context.Context) (*Application, error) {
	conf, err := config.New()
	if err != nil {
		return nil, err
	}

	app := &Application{}

	if err := app.setObserver(ctx, conf); err != nil {
		return nil, fmt.Errorf("set observers: %w", err)
	}

	if err := app.setStorages(ctx, conf); err != nil {
		return nil, fmt.Errorf("set storages: %w", err)
	}

	if err := app.setRepositories(conf); err != nil {
		return nil, fmt.Errorf("set repositories: %w", err)
	}

	if err := app.setServices(); err != nil {
		return nil, fmt.Errorf("set repositories: %w", err)
	}

	if err := app.setMiddlewares(); err != nil {
		return nil, fmt.Errorf("set middlewares: %w", err)
	}

	if err := app.setRouteHandlers(); err != nil {
		return nil, fmt.Errorf("set routehandlers: %w", err)
	}

	if err := app.setRoutes(conf); err != nil {
		return nil, err
	}

	if err := app.setServers(conf); err != nil {
		return nil, fmt.Errorf("set servers: %w", err)
	}

	return app, nil
}

func (a *Application) setObserver(ctx context.Context, conf *config.Configuration) error {
	l := logger.New(logger.Configuration{
		LogFormat:      logger.ParseLogFormat(conf.LogFormat),
		LogLevel:       logger.ParseLogLevel(conf.LogLevel),
		SkipFrameCount: 1,
		AppVersion:     conf.AppVersion,
		GitCommit:      conf.GitCommit,
	})

	a.observ = observer.New(ctx, observer.Configuration{Logger: l})

	return nil
}

func (a *Application) setStorages(ctx context.Context, conf *config.Configuration) error {
	maxLifetime := time.Duration(conf.PostgresConnMaxLifetimeMin) * time.Minute
	maxIdleTime := time.Duration(conf.PostgresConnMaxIdleTimeMin) * time.Minute
	db, err := database.New(ctx, database.Configuration{
		ConnectionURL:   conf.PostgresURL,
		MaxConnLifeTime: maxLifetime,
		MaxConnIdleTime: maxIdleTime,
	})

	if err != nil {
		return err
	}

	a.db = *db

	return nil
}

func (a *Application) setRepositories(conf *config.Configuration) error {
	a.txRepo = uow.New(a.db.Client)
	a.userRepository = user.New(a.db.Client)
	a.playerRepository = player.New(a.db.Client)
	a.teamRepository = team.New(a.db.Client)
	a.authRepository = auth.New(a.db.Client)
	a.transferRepository = transfer.New(a.db.Client)
	a.tokenRepository = token.New(token.Configuration{
		JWTExpiryMinutes: time.Duration(conf.JWTExpiryMinutes),
		JWTIssuer:        conf.JWTIssuer,
		JWTSecret:        conf.JWTSecret,
	})

	return nil
}

func (a *Application) setServices() error {
	localizationService, err := localization.New()
	if err != nil {
		return err
	}

	a.localizationService = localizationService
	a.txService = txService.New(txService.Configuration{
		TxManager: a.txRepo,
	})
	a.userService = userService.New(a.userRepository)
	a.tokenService = tokenService.New(a.tokenRepository)
	a.randomService = randomService.New()
	a.countryService = coutnryService.New()
	a.playerService = playerService.New(playerService.Configuration{
		TxService:      *a.txService,
		PlayerRepo:     a.playerRepository,
		RandomService:  *a.randomService,
		CountryService: *a.countryService,
	})
	a.teamService = teamService.New(teamService.Configuration{
		TxService:      *a.txService,
		TeamRepo:       a.teamRepository,
		UserService:    *a.userService,
		PlayerService:  *a.playerService,
		CountryService: *a.countryService,
	})
	a.transferService = transferService.New(transferService.Configuration{
		TransferRepo:  a.transferRepository,
		TxService:     *a.txService,
		PlayerService: *a.playerService,
	})
	a.authService = authService.New(authService.Configuration{
		TxService:    *a.txService,
		AuthRepo:     a.authRepository,
		UserService:  *a.userService,
		TokenService: *a.tokenService,
		TeamService:  *a.teamService,
	})

	return nil
}

func (a *Application) setMiddlewares() error {
	a.observMW = observMW.New(observMW.Configuration{
		Observer: a.observ,
	})
	a.localMW = localMW.New(localMW.Configuration{
		LocalizationService: a.localizationService,
	})
	a.authMW = authMW.New(authMW.Configuration{
		TokenService:        *a.tokenService,
		Observer:            a.observ,
		LocalizationService: a.localizationService,
	})

	return nil
}

func (a *Application) setRouteHandlers() error {
	a.authHandler = authHandler.New(authHandler.Configuration{
		AuthService:         a.authService,
		Observer:            a.observ,
		LocalizationService: a.localizationService,
	})

	a.teamHandler = teamHandler.New(teamHandler.Configuration{
		TeamService:         a.teamService,
		Observer:            a.observ,
		LocalizationService: a.localizationService,
	})

	a.playerHandler = playerHandler.New(playerHandler.Configuration{
		PlayerService:       a.playerService,
		Observer:            a.observ,
		LocalizationService: a.localizationService,
	})

	a.transferHandler = transferHandler.New(transferHandler.Configuration{
		TransferService:     a.transferService,
		Observer:            a.observ,
		LocalizationService: a.localizationService,
	})

	return nil
}

func (a *Application) setRoutes(conf *config.Configuration) error {
	r := chi.NewRouter()

	routeConfig := api.Configuration{
		Router:                 r,
		APIVersion:             conf.APIVersionPrefix,
		Observer:               a.observ,
		ObserverMiddleware:     a.observMW,
		AuthMiddleware:         a.authMW,
		LocalizationMiddleware: a.localMW,
		AuthHandler:            a.authHandler,
		TeamHandler:            a.teamHandler,
		TransferHandler:        a.transferHandler,
		PlayerHandler:          a.playerHandler,
		AuthService:            a.authService,
		TokenService:           a.tokenService,
		UserService:            a.userService,
		LocalizationService:    a.localizationService,
	}

	a.handler = api.New(routeConfig)

	return nil
}

func (a *Application) setServers(conf *config.Configuration) error {
	httpS, err := server.New(server.Configuration{
		Observer:     a.observ,
		Port:         conf.HTTPServerAddress,
		Handler:      a.handler,
		ReadTimeout:  conf.ReadTimeoutSeconds,
		WriteTimeout: conf.WriteTimeoutSeconds,
		IdleTimeout:  conf.IdleTimeoutSeconds,
	})

	if err != nil {
		return err
	}

	a.server = httpS
	return nil
}

func (a *Application) Start(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if err := a.server.Start(ctx); err != nil {
			return err
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func (a *Application) Close(ctx context.Context) error {
	if a.db.Client != nil {
		if err := a.db.Close(); err != nil {
			a.observ.Error(ctx, fmt.Errorf("closing database client: %w", err))
			return err
		}
	}

	if a.server != nil {
		if err := a.server.Close(ctx); err != nil {
			a.observ.Error(ctx, fmt.Errorf("closing http server: %w", err))
			return err
		}
	}

	return nil
}
