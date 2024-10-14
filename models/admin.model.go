// models/admin.model.go

package models

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"-"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type AdminResponse struct {
	ID        uuid.UUID    `json:"id,omitempty"`
	User      UserResponse `json:"user,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
}

type AssignAdminRoleRequest struct {
	UserID string `json:"userId" binding:"required"`
}
