package models

import (
	"soccer-api/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Password  string     `gorm:"not null"`
	Team      *TeamModel `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

func (um *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	if um.ID == uuid.Nil {
		um.ID = uuid.New()
	}
	return
}

func (um *UserModel) ToUserEntity() *entity.User {
	user := &entity.User{
		ID:        um.ID,
		Email:     um.Email,
		Password:  um.Password,
		CreatedAt: um.CreatedAt,
		UpdatedAt: um.UpdatedAt,
	}

	if um.Team != nil {
		user.Team = *um.Team.ToTeamEntity()
	}

	return user
}

func (um *UserModel) ToUserModal(e *entity.User) *UserModel {
	return &UserModel{
		ID:        e.ID,
		Email:     e.Email,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type TeamModel struct {
	ID        uuid.UUID     `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID     `gorm:"type:uuid;uniqueIndex;not null"`
	User      *UserModel    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name      *string       `gorm:"not null"`
	Country   *string       `gorm:"not null"`
	Budget    *int          `gorm:"not null"`
	Players   []PlayerModel `gorm:"foreignKey:TeamID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
}

func (TeamModel) TableName() string {
	return "teams"
}

func (tm *TeamModel) BeforeCreate(tx *gorm.DB) (err error) {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return
}

func (tm *TeamModel) ToTeamEntity() *entity.Team {
	teamEntity := &entity.Team{
		ID:        tm.ID,
		UserID:    tm.UserID,
		Name:      *tm.Name,
		Country:   entity.Country(*tm.Country),
		Budget:    *tm.Budget,
		CreatedAt: tm.CreatedAt,
		UpdatedAt: tm.UpdatedAt,
	}
	if len(tm.Players) > 0 {
		teamEntity.Players = make([]entity.Player, len(tm.Players))
		for i, playerModel := range tm.Players {
			teamEntity.Players[i] = *playerModel.ToPlayerEntity()
		}
	}
	return teamEntity
}

func (tm *TeamModel) ToTeamModal(e *entity.Team) *TeamModel {
	team := &TeamModel{
		ID:      e.ID,
		UserID:  e.UserID,
		Name:    &e.Name,
		Country: (*string)(&e.Country),
		Budget:  &e.Budget,
	}
	if e.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return team
}

type PlayerModel struct {
	ID          uuid.UUID             `gorm:"type:uuid;primary_key"`
	TeamID      uuid.UUID             `gorm:"type:uuid;uniqueIndex;not null"`
	Team        *TeamModel            `gorm:"foreignKey:TeamID;references:ID"`
	FirstName   string                `gorm:"not null"`
	LastName    string                `gorm:"not null"`
	Country     string                `gorm:"not null"`
	Age         int                   `gorm:"not null"`
	Position    entity.PlayerPosition `gorm:"type:varchar(2);not null"`
	MarketValue int                   `gorm:"not null"`
	CreatedAt   time.Time             `gorm:"autoCreateTime"`
	UpdatedAt   time.Time             `gorm:"autoUpdateTime"`
}

func (PlayerModel) TableName() string {
	return "players"
}

func (pm *PlayerModel) BeforeCreate(tx *gorm.DB) (err error) {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return
}

func (pm *PlayerModel) ToPlayerEntity() *entity.Player {
	return &entity.Player{
		ID:          pm.ID,
		TeamID:      pm.TeamID,
		FirstName:   pm.FirstName,
		LastName:    pm.LastName,
		Country:     pm.Country,
		Age:         pm.Age,
		Position:    pm.Position,
		MarketValue: pm.MarketValue,
		CreatedAt:   pm.CreatedAt,
		UpdatedAt:   pm.UpdatedAt,
	}
}

func (pm *PlayerModel) ToPlayerModal(e *entity.Player) *PlayerModel {
	player := &PlayerModel{
		ID:          e.ID,
		TeamID:      e.TeamID,
		FirstName:   e.FirstName,
		LastName:    e.LastName,
		Country:     e.Country,
		Age:         e.Age,
		Position:    e.Position,
		MarketValue: e.MarketValue,
	}
	if e.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return player
}

type TransferListingModel struct {
	ID            uuid.UUID    `gorm:"type:uuid;primary_key"`
	PlayerID      uuid.UUID    `gorm:"not null;unique"`
	Player        *PlayerModel `gorm:"foreignKey:PlayerID;references:ID"`
	SellingTeamID uuid.UUID    `gorm:"not null"`
	SellingTeam   *TeamModel   `gorm:"foreignKey:SellingTeamID;references:ID"`
	AskingPrice   int          `gorm:"not null"`
	ListedAt      time.Time    `gorm:"not null;autoCreateTime"`
	CreatedAt     time.Time    `gorm:"autoCreateTime"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime"`
}

func (TransferListingModel) TableName() string {
	return "transfer_listings"
}

func (tlm *TransferListingModel) BeforeCreate(tx *gorm.DB) (err error) {
	if tlm.ID == uuid.Nil {
		tlm.ID = uuid.New()
	}
	if tlm.ListedAt.IsZero() {
		tlm.ListedAt = time.Now()
	}
	return
}

func (tlm *TransferListingModel) ToListingEntity() *entity.TransferListing {
	listing := &entity.TransferListing{
		ID:            tlm.ID,
		PlayerID:      tlm.PlayerID,
		SellingTeamID: tlm.SellingTeamID,
		AskingPrice:   tlm.AskingPrice,
		ListedAt:      tlm.ListedAt,
		CreatedAt:     tlm.CreatedAt,
		UpdatedAt:     tlm.UpdatedAt,
	}
	if tlm.Player != nil {
		listing.Player = *tlm.Player.ToPlayerEntity()
	}
	if tlm.SellingTeam != nil {
		listing.SellingTeam = *tlm.SellingTeam.ToTeamEntity()
	}
	return listing
}

func (tlm *TransferListingModel) ToListingModel(e *entity.TransferListing) *TransferListingModel {
	id := e.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	model := &TransferListingModel{
		ID:            id,
		PlayerID:      e.PlayerID,
		SellingTeamID: e.SellingTeamID,
		AskingPrice:   e.AskingPrice,
		ListedAt:      e.ListedAt,
	}
	if model.ListedAt.IsZero() {
		model.ListedAt = time.Now()
	}
	return model
}
