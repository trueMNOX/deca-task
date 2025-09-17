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
type UserListResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int64          `json:"total"`
	Users []UserResponse `json:"users"`
}
