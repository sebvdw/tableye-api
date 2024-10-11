package models

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID        uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex" json:"user_id,omitempty"`
	User          User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Nickname      string        `gorm:"type:varchar(255);not null" json:"nickname,omitempty"`
	TotalWinnings float64       `gorm:"not null" json:"total_winnings,omitempty"`
	GamesPlayed   int           `gorm:"not null" json:"games_played,omitempty"`
	Rank          string        `gorm:"type:varchar(50)" json:"rank,omitempty"`
	Status        string        `gorm:"type:varchar(50);not null" json:"status,omitempty"`
	PlayedGames   []GameSummary `gorm:"many2many:game_players;" json:"played_games,omitempty"`
	WonGames      []GameSummary `gorm:"foreignKey:WinnerID" json:"won_games,omitempty"`
	CreatedAt     time.Time     `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time     `gorm:"not null" json:"updated_at,omitempty"`
}

type CreatePlayerRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type UpdatePlayerRequest struct {
	Nickname      string  `json:"nickname,omitempty"`
	TotalWinnings float64 `json:"total_winnings,omitempty"`
	GamesPlayed   int     `json:"games_played,omitempty"`
	Rank          string  `json:"rank,omitempty"`
	Status        string  `json:"status,omitempty"`
}

type PlayerResponse struct {
	ID            uuid.UUID    `json:"id,omitempty"`
	User          UserResponse `json:"user,omitempty"`
	Nickname      string       `json:"nickname,omitempty"`
	TotalWinnings float64      `json:"total_winnings,omitempty"`
	GamesPlayed   int          `json:"games_played,omitempty"`
	Rank          string       `json:"rank,omitempty"`
	Status        string       `json:"status,omitempty"`
	CreatedAt     time.Time    `json:"created_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at,omitempty"`
}
