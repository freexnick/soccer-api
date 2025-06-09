package response

import (
	"time"

	"github.com/google/uuid"

	"soccer-api/internal/domain/entity"
)

type PlayerSummaryForListing struct {
	ID          uuid.UUID             `json:"id"`
	FirstName   string                `json:"first_name"`
	LastName    string                `json:"last_name"`
	Country     entity.Country        `json:"country"`
	Age         int                   `json:"age"`
	Position    entity.PlayerPosition `json:"position"`
	MarketValue int                   `json:"market_value"`
}

type TeamSummaryForListing struct {
	ID      uuid.UUID      `json:"id"`
	Name    string         `json:"name"`
	Country entity.Country `json:"country"`
}

type TransferListItemResponse struct {
	ListingID   uuid.UUID               `json:"listing_id"`
	Player      PlayerSummaryForListing `json:"player"`
	SellingTeam TeamSummaryForListing   `json:"selling_team"`
	AskingPrice int                     `json:"asking_price"`
	ListedAt    time.Time               `json:"listed_at"`
}

type PlayerListedResponse struct {
	Message     string    `json:"message"`
	ListingID   uuid.UUID `json:"listing_id"`
	PlayerID    uuid.UUID `json:"player_id"`
	AskingPrice int       `json:"asking_price"`
}

type PlayerPurchasedResponse struct {
	Message string               `json:"message"`
	Player  PlayerDetailResponse `json:"player"`
}
