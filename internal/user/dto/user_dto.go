package dto

import "time"

type RequestUserDTO struct {
	Phone string `json:"phone" binding:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
