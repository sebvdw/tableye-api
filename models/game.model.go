package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name          string        `gorm:"uniqueIndex;not null" json:"name,omitempty"`
	Type          string        `gorm:"not null" json:"type,omitempty"`
	Description   string        `gorm:"type:text" json:"description,omitempty"`
	MaxPlayers    int           `gorm:"not null" json:"max_players,omitempty"`
	MinPlayers    int           `gorm:"not null" json:"min_players,omitempty"`
	MinBet        float64       `gorm:"not null" json:"min_bet,omitempty"`
	MaxBet        float64       `gorm:"not null" json:"max_bet,omitempty"`
	Casinos       []Casino      `gorm:"many2many:casino_games;" json:"casinos,omitempty"`
	GameSummaries []GameSummary `gorm:"foreignKey:GameID" json:"game_summaries,omitempty"`
	CreatedAt     time.Time     `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time     `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateGameRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description"`
	MaxPlayers  int     `json:"max_players" binding:"required"`
	MinPlayers  int     `json:"min_players" binding:"required"`
	MinBet      float64 `json:"min_bet" binding:"required"`
	MaxBet      float64 `json:"max_bet" binding:"required"`
}

type UpdateGameRequest struct {
	Name        string  `json:"name,omitempty"`
	Type        string  `json:"type,omitempty"`
	Description string  `json:"description,omitempty"`
	MaxPlayers  int     `json:"max_players,omitempty"`
	MinPlayers  int     `json:"min_players,omitempty"`
	MinBet      float64 `json:"min_bet,omitempty"`
	MaxBet      float64 `json:"max_bet,omitempty"`
}

type GameResponse struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
	MaxPlayers  int       `json:"max_players,omitempty"`
	MinPlayers  int       `json:"min_players,omitempty"`
	MinBet      float64   `json:"min_bet,omitempty"`
	MaxBet      float64   `json:"max_bet,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
