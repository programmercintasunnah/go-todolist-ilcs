// tests/task_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/programmercintasunnah/go-todolist-ilcs/controllers"
	"github.com/programmercintasunnah/go-todolist-ilcs/models"
	"github.com/programmercintasunnah/go-todolist-ilcs/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	logger := logrus.New()

	// Mock service
	mockService := &services.MockTaskService{}

	taskController := controllers.NewTaskController(mockService, logger)

	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		// Mock JWT authentication
		c.Next()
	})
	{
		protected.POST("/tasks", taskController.CreateTask)
		// Tambahkan endpoint lain jika perlu
	}

	return router
}

func TestCreateTask(t *testing.T) {
	router := SetupRouter()

	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Task created successfully", response["message"])
}
