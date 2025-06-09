package service

import (
	"context"
	"math"
	"math/rand"

	"github.com/google/uuid"

	"soccer-api/internal/domain/entity"
	"soccer-api/internal/domain/repository"
	playerService "soccer-api/internal/domain/service/player"
	txService "soccer-api/internal/domain/service/uow"
)

type Transfer struct {
	txService     txService.TxManager
	transferRepo  repository.Transfer
	playerService playerService.Player
}

func New(c Configuration) *Transfer {
	return &Transfer{
		transferRepo:  c.TransferRepo,
		txService:     c.TxService,
		playerService: c.PlayerService,
	}
}

func (t *Transfer) ListPlayer(ctx context.Context, userID, playerID uuid.UUID, askingPrice int) (*entity.TransferListing, error) {
	if askingPrice <= 0 {
		return nil, ErrPriceRequired
	}

	var listing *entity.TransferListing

	err := t.txService.Execute(ctx, func(repos repository.Repositories) error {
		player, err := repos.Player.GetByID(ctx, playerID)
		if err != nil {
			return err
		}
		if player == nil {
			return playerService.ErrPlayerNotFound
		}

		user, err := repos.User.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil || user.Team.ID != player.TeamID {
			return ErrListingPermissionDenied
		}

		existingListing, err := repos.Transfer.GetByPlayerID(ctx, playerID)
		if err != nil {
			return err
		}
		if existingListing != nil {
			return ErrPlayerAlreadyListed
		}

		newListing := &entity.TransferListing{
			ID:            uuid.New(),
			PlayerID:      playerID,
			SellingTeamID: player.TeamID,
			AskingPrice:   askingPrice,
		}

		if err := repos.Transfer.Create(ctx, newListing); err != nil {
			return err
		}

		listing = newListing
		listing.Player = *player
		listing.SellingTeam = user.Team

		return nil
	})

	if err != nil {
		return nil, err
	}

	return listing, nil
}

func (t *Transfer) ViewTransferList(ctx context.Context) ([]entity.TransferListing, error) {
	return t.transferRepo.GetAll(ctx)
}

func (t *Transfer) RemovePlayerFromTransfer(ctx context.Context, userID, playerID uuid.UUID) error {
	return t.txService.Execute(ctx, func(repos repository.Repositories) error {
		listing, err := repos.Transfer.GetByPlayerID(ctx, playerID)
		if err != nil {
			return err
		}
		if listing == nil {
			return ErrPlayerNotListed
		}

		user, err := repos.User.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil || user.Team.ID != listing.SellingTeamID {
			return ErrListingPermissionDenied
		}

		if err := repos.Transfer.Delete(ctx, listing.ID); err != nil {
			return err
		}

		return nil
	})
}

func (t *Transfer) BuyPlayer(ctx context.Context, buyerID, playerID uuid.UUID) (*entity.Player, error) {
	var purchasedPlayer *entity.Player
	uowErr := t.txService.Execute(ctx, func(repos repository.Repositories) error {
		buyer, err := repos.User.GetByID(ctx, buyerID)
		if err != nil {
			return err
		}

		listing, err := repos.Transfer.GetByPlayerID(ctx, playerID)
		if err != nil {
			return err
		}
		if listing == nil {
			return ErrPlayerNotListed
		}

		playerToBuy, err := repos.Player.GetByID(ctx, listing.PlayerID)
		if err != nil {
			return err
		}
		if playerToBuy == nil {
			return playerService.ErrPlayerNotFound
		}
		if listing.SellingTeamID == buyer.Team.ID {
			return ErrCannotBuyOwn
		}
		if buyer.Team.Budget < listing.AskingPrice {
			return ErrNotEnoughBudget
		}

		sellerTeam, err := repos.Team.GetTeamByID(ctx, listing.SellingTeamID)
		if err != nil {
			return err
		}
		sellerTeam.Budget += listing.AskingPrice
		if err := repos.Team.Update(ctx, sellerTeam); err != nil {
			return err
		}

		buyer.Team.Budget -= listing.AskingPrice
		if err := repos.Team.Update(ctx, &buyer.Team); err != nil {
			return err
		}

		priceBump := 0.10 + rand.Float64()*0.90
		playerToBuy.MarketValue = int(math.Round(float64(playerToBuy.MarketValue) * (1 + priceBump)))
		playerToBuy.TeamID = buyer.Team.ID
		if err := repos.Player.Update(ctx, playerToBuy); err != nil {
			return err
		}

		if err := repos.Transfer.Delete(ctx, listing.ID); err != nil {
			return err
		}

		purchasedPlayer = playerToBuy

		return nil
	})

	if uowErr != nil {
		return nil, uowErr
	}

	return purchasedPlayer, nil
}
