// services/mock_service.go
package services

import (
	"errors"

	"github.com/programmercintasunnah/go-todolist-ilcs/models"
	"github.com/programmercintasunnah/go-todolist-ilcs/repositories"
)

type MockTaskService struct {
	Tasks []models.Task
}

func (m *MockTaskService) CreateTask(task *models.Task) error {
	m.Tasks = append(m.Tasks, *task)
	return nil
}

func (m *MockTaskService) GetTaskByID(id uint) (*models.Task, error) {
	for _, task := range m.Tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, errors.New("task not found")
}

func (m *MockTaskService) GetAllTasks(filter map[string]interface{}, pagination repositories.Pagination, search string) ([]models.Task, int64, error) {
	return m.Tasks, int64(len(m.Tasks)), nil
}

func (m *MockTaskService) UpdateTask(id uint, updatedTask *models.Task) error {
	for i, task := range m.Tasks {
		if task.ID == id {
			m.Tasks[i] = *updatedTask
			return nil
		}
	}
	return errors.New("task not found")
}

func (m *MockTaskService) DeleteTask(id uint) error {
	for i, task := range m.Tasks {
		if task.ID == id {
			m.Tasks = append(m.Tasks[:i], m.Tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
