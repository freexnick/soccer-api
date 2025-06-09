package service

import (
	"context"

	"github.com/google/uuid"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
	countryService "soccer-api/internal/domain/service/country"
	randomService "soccer-api/internal/domain/service/random"
	txService "soccer-api/internal/domain/service/uow"
)

type Player struct {
	playerRepo     repository.Player
	txService      txService.TxManager
	randomService  randomService.Random
	countryService countryService.Country
}

func New(c Configuration) *Player {
	return &Player{
		playerRepo:     c.PlayerRepo,
		txService:      c.TxService,
		randomService:  c.RandomService,
		countryService: c.CountryService,
	}
}

func (p *Player) CreateInitialPlayers(ctx context.Context, teamID uuid.UUID) []entity.Player {
	players := make([]entity.Player, 0, entity.DefaultTeamPlayers)

	playerCounts := map[entity.PlayerPosition]int{
		entity.Goalkeeper: entity.InitialGoalkeepers,
		entity.Defender:   entity.InitialDefenders,
		entity.Midfielder: entity.InitialMidfielders,
		entity.Attacker:   entity.InitialAttackers,
	}

	for pos, count := range playerCounts {
		for range count {
			player := entity.Player{
				TeamID:      teamID,
				FirstName:   p.randomService.FirstName(ctx),
				LastName:    p.randomService.LastName(ctx),
				Country:     string(p.countryService.Random()),
				Age:         p.randomService.Age(ctx, entity.MinPlayerAge, entity.MaxPlayerAge),
				Position:    pos,
				MarketValue: entity.InitialPlayerValue,
			}
			players = append(players, player)
		}
	}

	return players
}

func (p *Player) UpdatePlayerDetails(ctx context.Context, userID, playerID uuid.UUID, firstName, lastName, country string) (*entity.Player, error) {
	if firstName == "" && lastName == "" && country == "" {
		return nil, ErrUpdateNoFields
	}

	var updatedPlayer *entity.Player

	uowErr := p.txService.Execute(ctx, func(repos repository.Repositories) error {
		player, err := repos.Player.GetByID(ctx, playerID)
		if err != nil {
			return err
		}
		if player == nil {
			return ErrPlayerNotFound
		}

		user, err := repos.User.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil || user.Team.ID != player.TeamID {
			return ErrPlayerNotOwned
		}

		if country != "" {
			parsedCountry, isValid := p.countryService.GetCountryByName(ctx, country)
			if !isValid {
				return ErrInvalidCountry
			}
			player.Country = string(parsedCountry)
		}

		if firstName != "" && player.FirstName != firstName {
			player.FirstName = firstName
		}
		if lastName != "" && player.LastName != lastName {
			player.LastName = lastName
		}

		if err := repos.Player.Update(ctx, player); err != nil {
			return err
		}

		updatedPlayer = player

		return nil
	})

	if uowErr != nil {
		return nil, uowErr
	}

	return updatedPlayer, nil
}

func (p *Player) GetPositionLabelKey(position entity.PlayerPosition) string {
	switch position {
	case entity.Goalkeeper:
		return LblPositionGoalkeeper
	case entity.Defender:
		return LblPositionDefender
	case entity.Midfielder:
		return LblPositionMidfielder
	case entity.Attacker:
		return LblPositionAttacker
	default:
		return string(position)
	}
}
