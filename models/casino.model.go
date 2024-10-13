package models

import (
	"time"

	"github.com/google/uuid"
)

type Casino struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name          string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"name,omitempty"`
	Location      string    `gorm:"type:varchar(255);not null" json:"location,omitempty"`
	LicenseNumber string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"license_number,omitempty"`
	Description   string    `gorm:"type:text" json:"description,omitempty"`
	OpeningHours  string    `gorm:"type:varchar(255)" json:"opening_hours,omitempty"`
	Website       string    `gorm:"type:varchar(255)" json:"website,omitempty"`
	PhoneNumber   string    `gorm:"type:varchar(50)" json:"phone_number,omitempty"`
	MaxCapacity   int       `gorm:"not null" json:"max_capacity,omitempty"`
	Status        string    `gorm:"type:varchar(50);not null" json:"status,omitempty"`
	Rating        float32   `json:"rating,omitempty"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateCasinoRequest struct {
	Name          string `json:"name" binding:"required"`
	Location      string `json:"location" binding:"required"`
	LicenseNumber string `json:"license_number" binding:"required"`
	Description   string `json:"description"`
	OpeningHours  string `json:"opening_hours"`
	Website       string `json:"website"`
	PhoneNumber   string `json:"phone_number"`
	MaxCapacity   int    `json:"max_capacity" binding:"required"`
	Status        string `json:"status" binding:"required"`
}

type UpdateCasinoRequest struct {
	Name          string  `json:"name,omitempty"`
	Location      string  `json:"location,omitempty"`
	LicenseNumber string  `json:"license_number,omitempty"`
	Description   string  `json:"description,omitempty"`
	OpeningHours  string  `json:"opening_hours,omitempty"`
	Website       string  `json:"website,omitempty"`
	PhoneNumber   string  `json:"phone_number,omitempty"`
	MaxCapacity   int     `json:"max_capacity,omitempty"`
	Status        string  `json:"status,omitempty"`
	Rating        float32 `json:"rating,omitempty"`
}

type CasinoResponse struct {
	ID            uuid.UUID `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Location      string    `json:"location,omitempty"`
	LicenseNumber string    `json:"license_number,omitempty"`
	Description   string    `json:"description,omitempty"`
	OpeningHours  string    `json:"opening_hours,omitempty"`
	Website       string    `json:"website,omitempty"`
	PhoneNumber   string    `json:"phone_number,omitempty"`
	MaxCapacity   int       `json:"max_capacity,omitempty"`
	Status        string    `json:"status,omitempty"`
	Rating        float32   `json:"rating,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
