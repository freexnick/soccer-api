package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransferListing struct {
	ID            uuid.UUID
	PlayerID      uuid.UUID
	Player        Player
	SellingTeamID uuid.UUID
	SellingTeam   Team
	AskingPrice   int
	ListedAt      time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
