package service

import (
	"github.com/KarmaBeLike/time-tracker-api/internal/external"
	"github.com/KarmaBeLike/time-tracker-api/internal/models"
	repositories "github.com/KarmaBeLike/time-tracker-api/internal/repository"
)

type UserService interface {
	GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error)
	CreateUser(passportNumber string) (*models.User, error)
	DeleteUser(userId string) error
	UpdateUser(user *models.User) error
}

type userService struct {
	userRepo        repositories.UserRepository
	peopleAPIClient *external.PeopleAPIClient
}

func NewUserService(userRepo repositories.UserRepository, peopleAPIClient *external.PeopleAPIClient) UserService {
	return &userService{
		userRepo:        userRepo,
		peopleAPIClient: peopleAPIClient,
	}
}

func (s *userService) CreateUser(passportNumber string) (*models.User, error) {
	peopleInfo, err := s.peopleAPIClient.GetPersonInfo(passportNumber)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		PassportNumber: passportNumber,
		Name:           peopleInfo.Name,
		Surname:        peopleInfo.Surname,
		Patronymic:     peopleInfo.Patronymic,
		Address:        peopleInfo.Address,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error) {
	return s.userRepo.GetUsers(page, limit, filters)
}

func (s *userService) DeleteUser(userId string) error {
	return s.userRepo.DeleteUser(userId)
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.UpdateUser(user)
}
