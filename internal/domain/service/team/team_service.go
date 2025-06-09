package service

import (
	"context"

	"github.com/google/uuid"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
	countryService "soccer-api/internal/domain/service/country"
	playerService "soccer-api/internal/domain/service/player"
	txService "soccer-api/internal/domain/service/uow"
	userService "soccer-api/internal/domain/service/user"
)

type Team struct {
	teamRepo       repository.Team
	txService      txService.TxManager
	userService    userService.User
	playerService  playerService.Player
	countryService countryService.Country
}

func New(c Configuration) *Team {
	return &Team{
		txService:      c.TxService,
		userService:    c.UserService,
		teamRepo:       c.TeamRepo,
		playerService:  c.PlayerService,
		countryService: c.CountryService,
	}
}

func (t *Team) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*entity.Team, int, error) {
	team, err := t.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, 0, err
	}
	if team == nil {
		return nil, 0, ErrUserHasNoTeam
	}

	totalValue := calculateTeamValue(team)

	return team, totalValue, nil
}

func (t *Team) GetTeamByUser(ctx context.Context, userID uuid.UUID) (*entity.Team, int, error) {
	user, err := t.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	team, err := t.teamRepo.GetTeamByUser(ctx, user.ID)
	if err != nil {
		return nil, 0, err
	}
	if team == nil {
		return nil, 0, ErrUserHasNoTeam
	}

	totalValue := calculateTeamValue(team)

	return team, totalValue, nil
}

func (t *Team) Create(ctx context.Context, userID uuid.UUID, teamName string) (*entity.Team, error) {
	team := &entity.Team{
		UserID:  userID,
		Name:    teamName,
		Country: t.countryService.Random(),
		Budget:  entity.InitialTeamBudget,
		Players: make([]entity.Player, 0, entity.DefaultTeamPlayers),
	}

	if err := t.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *Team) CreateTeamAndInitialPlayers(ctx context.Context, repos repository.Repositories, userID uuid.UUID, teamName string) (*entity.Team, error) {
	country := t.countryService.Random()

	team := &entity.Team{
		ID:      uuid.New(),
		UserID:  userID,
		Name:    teamName,
		Country: country,
		Budget:  entity.InitialTeamBudget,
	}

	if err := repos.Team.Create(ctx, team); err != nil {
		return nil, err
	}

	createdPlayers := t.playerService.CreateInitialPlayers(ctx, team.ID)

	if err := repos.Player.CreateBatch(ctx, createdPlayers, team.ID); err != nil {
		return nil, err
	}

	team.Players = createdPlayers
	return team, nil
}

func (t *Team) Update(ctx context.Context, userID uuid.UUID, name, country string) (*entity.Team, int, error) {
	if name == "" && country == "" {
		return nil, 0, ErrUpdateNoFields
	}

	var team *entity.Team

	uowErr := t.txService.Execute(ctx, func(repos repository.Repositories) error {
		user, err := repos.User.GetByID(ctx, userID)
		userTeam := user.Team
		if err != nil {
			return err
		}
		if user.Team.ID == uuid.Nil {
			return ErrUserHasNoTeam
		}

		if country != "" {
			parsedCountry, isValid := t.countryService.GetCountryByName(ctx, country)
			if !isValid {
				return ErrInvalidCountry
			}
			userTeam.Country = parsedCountry
		}

		userTeam.Name = name

		if err := repos.Team.Update(ctx, &userTeam); err != nil {
			return err
		}

		team = &userTeam

		return nil
	})

	if uowErr != nil {
		return nil, 0, uowErr
	}

	totalValue := calculateTeamValue(team)

	return team, totalValue, nil
}

func calculateTeamValue(team *entity.Team) int {
	if team == nil || len(team.Players) == 0 {
		return 0
	}
	var totalValue int
	for _, player := range team.Players {
		totalValue += player.MarketValue
	}
	return totalValue
}
