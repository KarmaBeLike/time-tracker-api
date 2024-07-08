package service

import (
	"time"

	"github.com/KarmaBeLike/time-tracker-api/internal/models"
	"github.com/KarmaBeLike/time-tracker-api/internal/repository"
)

type TaskService interface {
	GetWorklogs(userId int, startDate, endDate string) ([]models.Task, error)
	StartTask(userId, taskId int) error
	StopTask(userId, taskId int) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) GetWorklogs(userId int, startDate, endDate string) ([]models.Task, error) {
	return s.taskRepo.GetWorklogs(userId, startDate, endDate)
}

func (s *taskService) StartTask(userId, taskId int) error {
	return s.taskRepo.StartTask(userId, taskId, time.Now())
}

func (s *taskService) StopTask(userId, taskId int) error {
	return s.taskRepo.StopTask(userId, taskId, time.Now())
}
