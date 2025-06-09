package response

import (
	"time"

	"github.com/google/uuid"

	"soccer-api/internal/domain/entity"
)

type TeamResponse struct {
	ID         uuid.UUID             `json:"id"`
	UserID     uuid.UUID             `json:"user_id"`
	Name       string                `json:"name"`
	Country    entity.Country        `json:"country"`
	Budget     int                   `json:"budget"`
	TotalValue int                   `json:"total_value"`
	Players    []PlayerBasicResponse `json:"players"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}
