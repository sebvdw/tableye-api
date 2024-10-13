package models

import (
	"time"

	"github.com/google/uuid"
)

type GameSummary struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	GameID       uuid.UUID     `gorm:"type:uuid;not null" json:"-"`
	CasinoID     uuid.UUID     `gorm:"type:uuid;not null" json:"-"`
	DealerID     uuid.UUID     `gorm:"type:uuid;not null" json:"-"`
	Dealer       Dealer        `gorm:"foreignKey:DealerID" json:"dealer,omitempty"`
	Players      []Player      `gorm:"many2many:game_players;" json:"players,omitempty"`
	StartTime    time.Time     `gorm:"not null" json:"start_time,omitempty"`
	EndTime      time.Time     `json:"end_time,omitempty"`
	TotalPot     float64       `json:"total_pot,omitempty"`
	Status       string        `gorm:"type:varchar(50);not null" json:"status,omitempty"`
	RoundsPlayed int           `json:"rounds_played,omitempty"`
	HighestBet   float64       `json:"highest_bet,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:GameSummaryID" json:"transactions,omitempty"`
	CreatedAt    time.Time     `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time     `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateGameSummaryRequest struct {
	GameID    string    `json:"game_id" binding:"required"`
	CasinoID  string    `json:"casino_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	PlayerIDs []string  `json:"player_ids" binding:"required"`
	DealerID  string    `json:"dealer_id" binding:"required"`
}

type UpdateGameSummaryRequest struct {
	EndTime      time.Time `json:"end_time,omitempty"`
	TotalPot     float64   `json:"total_pot,omitempty"`
	Status       string    `json:"status,omitempty"`
	RoundsPlayed int       `json:"rounds_played,omitempty"`
	HighestBet   float64   `json:"highest_bet,omitempty"`
}

type GameSummaryResponse struct {
	ID           uuid.UUID        `json:"id,omitempty"`
	Game         GameResponse     `json:"game,omitempty"`
	Casino       CasinoResponse   `json:"casino,omitempty"`
	StartTime    time.Time        `json:"start_time,omitempty"`
	EndTime      time.Time        `json:"end_time,omitempty"`
	Players      []PlayerResponse `json:"players,omitempty"`
	Dealer       DealerResponse   `json:"dealer,omitempty"`
	TotalPot     float64          `json:"total_pot,omitempty"`
	Status       string           `json:"status,omitempty"`
	RoundsPlayed int              `json:"rounds_played,omitempty"`
	HighestBet   float64          `json:"highest_bet,omitempty"`
	CreatedAt    time.Time        `json:"created_at,omitempty"`
	UpdatedAt    time.Time        `json:"updated_at,omitempty"`
}
