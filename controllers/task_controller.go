// controllers/task_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/programmercintasunnah/go-todolist-ilcs/models"
	"github.com/programmercintasunnah/go-todolist-ilcs/repositories"
	"github.com/programmercintasunnah/go-todolist-ilcs/services"
	"github.com/sirupsen/logrus"
)

type TaskController struct {
	service services.TaskService
	logger  *logrus.Logger
}

func NewTaskController(service services.TaskService, logger *logrus.Logger) *TaskController {
	return &TaskController{
		service: service,
		logger:  logger,
	}
}

type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=pending completed"`
	DueDate     string `json:"due_date" binding:"required,datetime=2006-01-02"`
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var input CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		tc.logger.Error("CreateTask: Invalid input", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dueDate, err := time.Parse("2006-01-02", input.DueDate)
	if err != nil {
		tc.logger.Error("CreateTask: Invalid due_date format", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format"})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		DueDate:     dueDate,
	}

	if err := tc.service.CreateTask(&task); err != nil {
		tc.logger.Error("CreateTask: Failed to create task", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    task,
	})
}

type GetAllTasksQuery struct {
	Status string `form:"status"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	var query GetAllTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		tc.logger.Error("GetAllTasks: Invalid query parameters", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default pagination
	page := query.Page
	if page == 0 {
		page = 1
	}
	limit := query.Limit
	if limit == 0 {
		limit = 10
	}

	filter := make(map[string]interface{})
	if query.Status != "" {
		filter["status"] = query.Status
	}

	tasks, total, err := tc.service.GetAllTasks(filter, repositories.Pagination{Page: page, Limit: limit}, query.Search)
	if err != nil {
		tc.logger.Error("GetAllTasks: Failed to retrieve tasks", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_tasks":  total,
		},
	})
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		tc.logger.Error("GetTaskByID: Invalid ID", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := tc.service.GetTaskByID(uint(id))
	if err != nil {
		tc.logger.Error("GetTaskByID: Task not found", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=pending completed"`
	DueDate     string `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		tc.logger.Error("UpdateTask: Invalid ID", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var input UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		tc.logger.Error("UpdateTask: Invalid input", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dueDate time.Time
	if input.DueDate != "" {
		dueDate, err = time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			tc.logger.Error("UpdateTask: Invalid due_date format", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format"})
			return
		}
	}

	updatedTask := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		DueDate:     dueDate,
	}

	if err := tc.service.UpdateTask(uint(id), &updatedTask); err != nil {
		tc.logger.Error("UpdateTask: Failed to update task", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	task, _ := tc.service.GetTaskByID(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    task,
	})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		tc.logger.Error("DeleteTask: Invalid ID", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := tc.service.DeleteTask(uint(id)); err != nil {
		tc.logger.Error("DeleteTask: Failed to delete task", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
