package entity

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	UserID    uuid.UUID
	Email     string
	Issuer    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}
