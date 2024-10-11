package models

import (
	"time"

	"github.com/google/uuid"
)

type Dealer struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID        uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex" json:"user_id,omitempty"`
	User          User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DealerCode    string        `gorm:"uniqueIndex;not null" json:"dealer_code,omitempty"`
	Status        string        `gorm:"type:varchar(255);not null" json:"status,omitempty"`
	GamesDealt    int           `gorm:"not null" json:"games_dealt,omitempty"`
	Rating        float32       `gorm:"not null" json:"rating,omitempty"`
	Casinos       []Casino      `gorm:"many2many:casino_dealers;" json:"casinos,omitempty"`
	GameSummaries []GameSummary `gorm:"foreignKey:DealerID" json:"game_summaries,omitempty"`
	LastActiveAt  time.Time     `json:"last_active_at,omitempty"`
	CreatedAt     time.Time     `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time     `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateDealerRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	DealerCode string `json:"dealer_code" binding:"required"`
	Status     string `json:"status" binding:"required"`
}

type UpdateDealerRequest struct {
	Status       string    `json:"status,omitempty"`
	GamesDealt   int       `json:"games_dealt,omitempty"`
	Rating       float32   `json:"rating,omitempty"`
	LastActiveAt time.Time `json:"last_active_at,omitempty"`
}

type DealerResponse struct {
	ID           uuid.UUID    `json:"id,omitempty"`
	User         UserResponse `json:"user,omitempty"`
	DealerCode   string       `json:"dealer_code,omitempty"`
	Status       string       `json:"status,omitempty"`
	GamesDealt   int          `json:"games_dealt,omitempty"`
	Rating       float32      `json:"rating,omitempty"`
	LastActiveAt time.Time    `json:"last_active_at,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
	UpdatedAt    time.Time    `json:"updated_at,omitempty"`
}
