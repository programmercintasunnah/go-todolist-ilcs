// cmd/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/programmercintasunnah/go-todolist-ilcs/config"
	"github.com/programmercintasunnah/go-todolist-ilcs/controllers"
	"github.com/programmercintasunnah/go-todolist-ilcs/middlewares"
	"github.com/programmercintasunnah/go-todolist-ilcs/repositories"
	"github.com/programmercintasunnah/go-todolist-ilcs/services"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"         // PostgreSQL driver
	_ "github.com/godror/godror"  // Oracle driver
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
	
	cfg := config.LoadConfig()

	// Setup Logging
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Setup Database Connection
	var dsn string
	switch os.Getenv("DB_TYPE") {
	case "postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	case "oracle":
		dsn = fmt.Sprintf("oracle://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	default:
		log.Fatal("Unsupported DB_TYPE in .env file")
	}

	db, err := sql.Open(os.Getenv("DB_TYPE"), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize Repositories, Services, Controllers
	redisClient := repositories.InitRedis(cfg)
	taskRepo := repositories.NewTaskRepository(db, redisClient, logger)
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
