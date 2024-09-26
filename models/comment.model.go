package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateCommentRequest struct {
	Content string    `json:"content" binding:"required"`
	PostID  uuid.UUID `json:"post_id" binding:"required"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
}

type UpdateComment struct {
	Content string `json:"content,omitempty"`
}
