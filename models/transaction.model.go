package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	GameSummaryID uuid.UUID `gorm:"type:uuid;not null" json:"id,omitempty"`
	PlayerID      uuid.UUID `gorm:"type:uuid;not null" json:"id,omitempty"`
	Player        Player    `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount,omitempty"`
	Type          string    `gorm:"type:varchar(50);not null" json:"-"`
	Outcome       string    `gorm:"type:varchar(10);not null;default:'loss'" json:"outcome,omitempty"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateTransactionRequest struct {
	GameSummaryID string  `json:"game_summary_id" binding:"required"`
	PlayerID      string  `json:"player_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	//Type          string  `json:"type" binding:"required"`
	Outcome string `json:"outcome" binding:"required,oneof=win loss"`
}

type UpdateTransactionRequest struct {
	Amount float64 `json:"amount,omitempty"`
	//Type    string  `json:"type,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}

type TransactionResponse struct {
	ID        uuid.UUID      `json:"id"`
	Player    PlayerResponse `json:"player"`
	Amount    float64        `json:"amount"`
	Outcome   string         `json:"outcome"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
