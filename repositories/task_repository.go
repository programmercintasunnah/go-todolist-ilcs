package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/programmercintasunnah/go-todolist-ilcs/models"
	"github.com/sirupsen/logrus"
)

type Pagination struct {
	Page  int
	Limit int
}

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id uint) (*models.Task, error)
	GetAllTasks(filter map[string]interface{}, pagination Pagination, search string) ([]models.Task, int64, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id uint) error
}

type taskRepository struct {
	db          *sql.DB
	redisClient *redis.Client
	logger      *logrus.Logger
}

func NewTaskRepository(db *sql.DB, redisClient *redis.Client, logger *logrus.Logger) TaskRepository {
	return &taskRepository{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	query := "INSERT INTO tasks (title, description, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, task.Title, task.Description, task.Status, time.Now(), time.Now())
	return err
}

func (r *taskRepository) GetTaskByID(id uint) (*models.Task, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("task:%d", id)

	// Cek Cache
	cachedTask, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var task models.Task
		if err := json.Unmarshal([]byte(cachedTask), &task); err == nil {
			return &task, nil
		}
	}

	// Ambil dari Database
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var task models.Task
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	// Simpan ke Cache
	taskJSON, _ := json.Marshal(task)
	r.redisClient.Set(ctx, cacheKey, taskJSON, time.Minute*10)

	return &task, nil
}

func (r *taskRepository) GetAllTasks(filter map[string]interface{}, pagination Pagination, search string) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	// Base Query
	query := "SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM tasks WHERE 1=1"

	var args []interface{}

	// Filter Berdasarkan Status
	if status, ok := filter["status"]; ok {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, status)
	}

	// Search Berdasarkan Title atau Description
	if search != "" {
		query += " AND (title LIKE ? OR description LIKE ?)"
		countQuery += " AND (title LIKE ? OR description LIKE ?)"
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam)
	}

	// Pagination
	offset := (pagination.Page - 1) * pagination.Limit
	query += " LIMIT ? OFFSET ?"
	args = append(args, pagination.Limit, offset)

	// Eksekusi Count Query
	row := r.db.QueryRow(countQuery, args...)
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	// Eksekusi Query untuk Data
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func (r *taskRepository) UpdateTask(task *models.Task) error {
	query := "UPDATE tasks SET title = ?, description = ?, status = ?, updated_at = ? WHERE id = ?"
	_, err := r.db.Exec(query, task.Title, task.Description, task.Status, time.Now(), task.ID)
	return err
}

func (r *taskRepository) DeleteTask(id uint) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
