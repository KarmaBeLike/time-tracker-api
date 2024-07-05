package service

import (
	"github.com/KarmaBeLike/time-tracker-api/internal/models"
	repositories "github.com/KarmaBeLike/time-tracker-api/internal/repository"
)

type UserService interface {
	GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error) {
	return s.userRepo.GetUsers(page, limit, filters)
}
