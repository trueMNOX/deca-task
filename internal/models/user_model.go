package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PhoneNumber  string         `gorm:"unique;not null" json:"phone_number"`
	RegisteredAt time.Time      `gorm:"autoCreateTime" json:"registered_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Role string `json:"role"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
