package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlayerPosition string

const (
	Goalkeeper         PlayerPosition = "GK"
	Defender           PlayerPosition = "DF"
	Midfielder         PlayerPosition = "MF"
	Attacker           PlayerPosition = "AT"
	MinPlayerAge       int            = 18
	MaxPlayerAge       int            = 40
	InitialPlayerValue int            = 1_000_000
)

type Player struct {
	ID          uuid.UUID
	TeamID      uuid.UUID
	FirstName   string
	LastName    string
	Country     string
	Age         int
	Position    PlayerPosition
	MarketValue int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
