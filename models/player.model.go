package models

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Nickname      string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"nickname,omitempty"`
	TotalWinnings float64   `gorm:"not null" json:"total_winnings,omitempty"`
	Rank          string    `gorm:"type:varchar(50)" json:"rank,omitempty"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status,omitempty"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreatePlayerRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

type UpdatePlayerRequest struct {
	Nickname      string  `json:"nickname,omitempty"`
	TotalWinnings float64 `json:"total_winnings,omitempty"`
	Rank          string  `json:"rank,omitempty"`
	Status        string  `json:"status,omitempty"`
}

type PlayerResponse struct {
	ID            uuid.UUID `json:"id,omitempty"`
	Nickname      string    `json:"nickname,omitempty"`
	TotalWinnings float64   `json:"total_winnings,omitempty"`
	Rank          string    `json:"rank,omitempty"`
	Status        string    `json:"status,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
