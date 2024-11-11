// services/task_service.go
package services

import (
	"strings"

	"github.com/programmercintasunnah/go-todolist-ilcs/models"
	"github.com/programmercintasunnah/go-todolist-ilcs/repositories"
	"github.com/sirupsen/logrus"
)

type TaskService interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id uint) (*models.Task, error)
	GetAllTasks(filter map[string]interface{}, pagination repositories.Pagination, search string) ([]models.Task, int64, error)
	UpdateTask(id uint, updatedTask *models.Task) error
	DeleteTask(id uint) error
}

type taskService struct {
	repo   repositories.TaskRepository
	logger *logrus.Logger
}

func NewTaskService(repo repositories.TaskRepository, logger *logrus.Logger) TaskService {
	return &taskService{
		repo:   repo,
		logger: logger,
	}
}

func (s *taskService) CreateTask(task *models.Task) error {
	return s.repo.CreateTask(task)
}

func (s *taskService) GetTaskByID(id uint) (*models.Task, error) {
	return s.repo.GetTaskByID(id)
}

func (s *taskService) GetAllTasks(filter map[string]interface{}, pagination repositories.Pagination, search string) ([]models.Task, int64, error) {
	return s.repo.GetAllTasks(filter, pagination, search)
}

func (s *taskService) UpdateTask(id uint, updatedTask *models.Task) error {
	existingTask, err := s.repo.GetTaskByID(id)
	if err != nil {
		return err
	}

	if strings.TrimSpace(updatedTask.Title) != "" {
		existingTask.Title = updatedTask.Title
	}
	if strings.TrimSpace(updatedTask.Description) != "" {
		existingTask.Description = updatedTask.Description
	}
	if strings.TrimSpace(updatedTask.Status) != "" {
		existingTask.Status = updatedTask.Status
	}
	if !updatedTask.DueDate.IsZero() {
		existingTask.DueDate = updatedTask.DueDate
	}

	return s.repo.UpdateTask(existingTask)
}

func (s *taskService) DeleteTask(id uint) error {
	return s.repo.DeleteTask(id)
}
