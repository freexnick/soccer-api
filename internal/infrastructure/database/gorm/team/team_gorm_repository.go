package team

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/infrastructure/database/gorm/models"
)

type Team struct {
	model models.TeamModel
	db    *gorm.DB
}

func New(db *gorm.DB) *Team {
	return &Team{db: db}
}

func (t *Team) Create(ctx context.Context, team *entity.Team) error {
	teamModel := t.model.ToTeamModal(team)
	err := t.db.WithContext(ctx).Create(teamModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (t *Team) Delete(ctx context.Context, teamID uuid.UUID) error {
	return t.db.WithContext(ctx).Delete(&models.TeamModel{}, teamID).Error
}

func (t *Team) Update(ctx context.Context, team *entity.Team) error {
	teamModel := t.model.ToTeamModal(team)
	if err := t.db.WithContext(ctx).
		Model(&models.TeamModel{}).
		Where("id = ?", team.ID).
		Updates(teamModel).Error; err != nil {
		return err
	}

	return nil
}

func (t *Team) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*entity.Team, error) {
	var teamModel models.TeamModel

	err := t.db.WithContext(ctx).
		Preload("Players").
		First(&teamModel, "id = ?", teamID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return teamModel.ToTeamEntity(), nil
}

func (t *Team) GetTeamByUser(ctx context.Context, userID uuid.UUID) (*entity.Team, error) {
	var teamModel models.TeamModel

	err := t.db.WithContext(ctx).
		Preload("Players").
		Where("user_id = ?", userID).
		First(&teamModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return teamModel.ToTeamEntity(), nil
}
