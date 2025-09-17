package user

import "deca-task/internal/models"

type userService struct {
	userRepository *userRepository
}

func NewUserService(userRepository *userRepository) *userService {
	return &userService{userRepository: userRepository}
}

func (s *userService) FindUserById(userid uint) (*models.User, error) {
	user, err := s.userRepository.FindUserById(userid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *userService) FindUsers(page, limit int) ([]models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	users, total, err := s.userRepository.GetUsers(page, limit)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
