package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/programmercintasunnah/go-todolist-ilcs/config"
	"github.com/programmercintasunnah/go-todolist-ilcs/controllers"
	"github.com/programmercintasunnah/go-todolist-ilcs/middlewares"
	"github.com/programmercintasunnah/go-todolist-ilcs/repositories"
	"github.com/programmercintasunnah/go-todolist-ilcs/services"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	// Setup Logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Setup Database Connection
	dsn := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBService)
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Oracle database: %v", err)
	}
	defer db.Close()

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping Oracle database: %v", err)
	}

	// Initialize Repositories, Services, Controllers
	redisClient := repositories.InitRedis(cfg)
	taskRepo := repositories.NewTaskRepository(db, redisClient, logger) // Update repository to accept *sql.DB
	taskService := services.NewTaskService(taskRepo, logger)
	taskController := controllers.NewTaskController(taskService, logger)

	// Setup Gin
	router := gin.New()

	// Apply Middlewares
	router.Use(gin.Recovery())
	router.Use(middlewares.Logger(logger))

	// Public Routes
	public := router.Group("/api")
	{
		public.POST("/login", controllers.Login)
	}

	// Protected Routes
	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuth(cfg.JWTSecret))
	{
		protected.POST("/tasks", taskController.CreateTask)
		protected.GET("/tasks", taskController.GetAllTasks)
		protected.GET("/tasks/:id", taskController.GetTaskByID)
		protected.PUT("/tasks/:id", taskController.UpdateTask)
		protected.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	// Start Server
	if err := router.Run(":8080"); err != nil {
		logger.Fatal("Failed to run server: ", err)
	}
}
