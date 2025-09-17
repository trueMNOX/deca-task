package auth

import (
	"deca-task/internal/auth/dto"
	"deca-task/internal/auth/jwt"
	"fmt"
	"strconv"
	"time"
)

type authservice struct {
	repo *authrepository
}

func NewAuthService(repo *authrepository) *authservice {
	return &authservice{repo: repo}
}

func (s *authservice) LoginUser(phone uint) (*dto.RequestOTPResponse, error) {
	if err := s.repo.IncrementOtpRequestCount(phone, 10*time.Minute, 3); err != nil {
		if err == ErrRateLimitExceeded {
			return nil, fmt.Errorf("too many OTP requests, try again later")
		}
		return nil, err
	}
	Otp, err := s.repo.SaveOTP(phone)
	if err != nil {
		return nil, err
	}

	return &dto.RequestOTPResponse{
		Message: "OTP sent successfully",
		OTP:     Otp,
	}, nil
}
func (s *authservice) VerifyOTP(phone uint, otp string) (*dto.VerifyOTPResponse, error) {
	code, err := s.repo.GetOtpFromRedis(strconv.FormatUint(uint64(phone), 10))
	if err != nil {
		return nil, err
	}
	if code != otp {
		return nil, fmt.Errorf("invalid OTP")
	}
	user, err := s.repo.FindUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err := s.repo.SaveUser(phone)
		if err != nil {
			return nil, err
		}
		token, err := jwt.GenerateToken(user.ID)
		if err != nil {
			return nil, err
		}
		return &dto.VerifyOTPResponse{Token: token}, nil
	}
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.VerifyOTPResponse{Token: token}, nil

}
