package response

import (
	"time"

	"github.com/google/uuid"

	"soccer-api/internal/domain/entity"
)

type PlayerBasicResponse struct {
	ID          uuid.UUID             `json:"id"`
	FirstName   string                `json:"first_name"`
	LastName    string                `json:"last_name"`
	Country     entity.Country        `json:"country"`
	Age         int                   `json:"age"`
	Position    entity.PlayerPosition `json:"position"`
	MarketValue int                   `json:"market_value"`
}

type PlayerDetailResponse struct {
	ID          uuid.UUID             `json:"id"`
	TeamID      uuid.UUID             `json:"team_id,omitempty"`
	FirstName   string                `json:"first_name"`
	LastName    string                `json:"last_name"`
	Age         int                   `json:"age"`
	Position    entity.PlayerPosition `json:"position"`
	MarketValue int                   `json:"market_value"`
	Country     entity.Country        `json:"country"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}
