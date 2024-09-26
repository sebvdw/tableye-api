package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email,omitempty"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"type:varchar(255);not null" json:"role,omitempty"`
	Provider  string    `gorm:"not null" json:"provider,omitempty"`
	Verified  bool      `gorm:"not null" json:"verified,omitempty"`
	Posts     []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments  []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Provider  string    `json:"provider,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
