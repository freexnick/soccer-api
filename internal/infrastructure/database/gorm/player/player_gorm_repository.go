package player

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/infrastructure/database/gorm/models"
)

type Player struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Player {
	return &Player{db: db}
}

func (p *Player) Create(ctx context.Context, player *entity.Player) error {
	var playerModel models.PlayerModel

	if err := p.db.WithContext(ctx).Create(&playerModel).Error; err != nil {
		return err
	}

	return nil
}

func (p *Player) CreateBatch(ctx context.Context, players []entity.Player, teamID uuid.UUID) error {
	playerModels := make([]models.PlayerModel, len(players))
	for i, player := range players {
		var p models.PlayerModel
		playerModel := p.ToPlayerModal(&player)
		playerModels[i] = *playerModel
	}

	if err := p.db.WithContext(ctx).CreateInBatches(playerModels, len(playerModels)).Error; err != nil {
		return err
	}
	return nil
}

func (p *Player) Update(ctx context.Context, player *entity.Player) error {
	var playerModel models.PlayerModel

	if err := p.db.WithContext(ctx).
		Model(&models.PlayerModel{}).
		Where("id = ? ", player.ID).
		Updates(playerModel.ToPlayerModal(player)).Error; err != nil {
		return err
	}

	return nil
}

func (p *Player) Delete(ctx context.Context, playerID uuid.UUID) error {
	return nil
}

func (p *Player) GetByID(ctx context.Context, playerID uuid.UUID) (*entity.Player, error) {
	var playerModel models.PlayerModel

	if err := p.db.WithContext(ctx).
		Where("id = ?", playerID).
		First(&playerModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return playerModel.ToPlayerEntity(), nil
}
