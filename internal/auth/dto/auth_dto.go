package dto

import "time"

type RequestOTPDTO struct {
	Phone uint `json:"phone" binding:"required"`
}

type RequestOTPResponse struct {
	Message string `json:"message"`
	OTP   string `json:"otp"`
}

type VerifyOTPDTO struct {
	Phone uint `json:"phone" binding:"required"`
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
