package transfer

import (
	"context"
	"errors"
	"soccer-api/internal/domain/entity"
	"soccer-api/internal/infrastructure/database/gorm/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transfer struct {
	db    *gorm.DB
	model models.TransferListingModel
}

func New(db *gorm.DB) *Transfer {
	return &Transfer{db: db}
}

func (t *Transfer) Create(ctx context.Context, listing *entity.TransferListing) error {
	listingModel := t.model.ToListingModel(listing)

	if err := t.db.WithContext(ctx).Create(&listingModel).Error; err != nil {
		return err
	}

	return nil
}

func (t *Transfer) GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*entity.TransferListing, error) {
	var listingModel models.TransferListingModel

	if err := t.db.WithContext(ctx).
		Preload("Player").
		Preload("SellingTeam").
		Where("player_id = ?", playerID).
		First(&listingModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return listingModel.ToListingEntity(), nil
}

func (t *Transfer) Delete(ctx context.Context, listingID uuid.UUID) error {
	if err := t.db.WithContext(ctx).Delete(&models.TransferListingModel{}, listingID).Error; err != nil {
		return err
	}

	return nil
}

func (t *Transfer) Update(ctx context.Context, listing *entity.TransferListing) error {
	return nil
}

func (t *Transfer) GetAll(ctx context.Context) ([]entity.TransferListing, error) {
	var listingModels []models.TransferListingModel

	if err := t.db.WithContext(ctx).
		Preload("Player.Team").
		Preload("SellingTeam").
		Order("asking_price asc, listed_at asc").
		Find(&listingModels).Error; err != nil {
		return nil, err
	}

	listings := make([]entity.TransferListing, len(listingModels))
	for i, lm := range listingModels {
		listings[i] = *lm.ToListingEntity()
	}

	return listings, nil
}

func (t *Transfer) Delet(ctx context.Context, listingID uuid.UUID) error {
	if err := t.db.WithContext(ctx).Delete(&models.TransferListingModel{}, listingID).Error; err != nil {
		return err
	}

	return nil
}
