package user

import (
	"deca-task/internal/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUserById(userid uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUsers(page, limit int, phone string) ([]models.User, int, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})
	if phone != "" {
		query = query.Where("phone_number LIKE ?", "%"+phone+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	result := r.db.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return users, int(total), nil
}
