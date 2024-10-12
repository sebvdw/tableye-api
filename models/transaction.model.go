package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	GameSummaryID uuid.UUID   `gorm:"type:uuid;not null" json:"game_summary_id,omitempty"`
	GameSummary   GameSummary `gorm:"foreignKey:GameSummaryID" json:"game_summary,omitempty"`
	PlayerID      uuid.UUID   `gorm:"type:uuid;not null" json:"player_id,omitempty"`
	Player        Player      `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Amount        float64     `gorm:"not null" json:"amount,omitempty"`
	Type          string      `gorm:"type:varchar(50);not null" json:"type,omitempty"` // "bet" or "win"
	CreatedAt     time.Time   `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time   `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateTransactionRequest struct {
	GameSummaryID string  `json:"game_summary_id" binding:"required"`
	PlayerID      string  `json:"player_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Type          string  `json:"type" binding:"required"`
}

type UpdateTransactionRequest struct {
	Amount float64 `json:"amount,omitempty"`
	Type   string  `json:"type,omitempty"`
}

type TransactionResponse struct {
	ID            uuid.UUID `json:"id,omitempty"`
	GameSummaryID uuid.UUID `json:"game_summary_id,omitempty"`
	PlayerID      uuid.UUID `json:"player_id,omitempty"`
	Amount        float64   `json:"amount,omitempty"`
	Type          string    `json:"type,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
