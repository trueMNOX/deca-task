package dto

import "time"

type RequestOTPDTO struct {
	Phone string `json:"phone" binding:"required"`
}

type RequestOTPResponse struct {
	Message string `json:"message"`
}

type VerifyOTPDTO struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
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
