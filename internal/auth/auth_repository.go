package auth

import (
	"context"
	"deca-task/internal/models"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type authrepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAuthRepository(db *gorm.DB, redis *redis.Client) *authrepository {
	return &authrepository{
		db:    db,
		redis: redis,
	}
}

func (r *authrepository) FindUserByPhone(phone uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("phone_number = ?", strconv.FormatUint(uint64(phone), 10)).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authrepository) SaveUser(phone uint) (*models.User, error) {
	user := &models.User{
		PhoneNumber: strconv.FormatUint(uint64(phone), 10),
	}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authrepository) GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (r *authrepository) SaveOTP(phone uint) (string, error) {
	otp := r.GenerateOTP()
	ctx := context.Background()

	key := fmt.Sprintf("otp:%d", phone)
	err := r.redis.Set(ctx, key, otp, 2*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("failed to save OTP: %v", err)
	}

	return otp, nil
}
func (r *authrepository) GetOtpFromRedis(phone string) (string, error) {
	ctx := context.Background()

	key := fmt.Sprintf("otp:%s", phone)
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("OTP not found for phone %s", phone)
		}
		return "", fmt.Errorf("failed to get OTP: %v", err)
	}
	return val, nil
}
