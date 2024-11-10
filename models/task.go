// models/task.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" validate:"required"`
	Description string         `json:"description"`
	Status      string         `json:"status" validate:"oneof=pending completed"`
	DueDate     time.Time      `json:"due_date" validate:"required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
