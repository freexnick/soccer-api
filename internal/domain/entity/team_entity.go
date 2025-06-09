package entity

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Country   Country
	Budget    int
	Players   []Player
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	InitialTeamBudget int = 5_000_000

	DefaultTeamPlayers = 20
	InitialGoalkeepers = 3
	InitialDefenders   = 6
	InitialMidfielders = 6
	InitialAttackers   = 5
)
