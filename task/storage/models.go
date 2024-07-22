package storage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TaskStatusCreated    = "created"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
)

type Task struct {
	ID          string    `sql:"size:255;not null";gorm:"primary_key"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:50;not null;default:'pending'" json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UserID      string    `gorm:"type:string;index:idx_task_user_id;not null"`
}

func (task *Task) BeforeSave(tx *gorm.DB) (err error) {
	if task.ID == "" {
		task.ID = uuid.NewString()
	}
	if task.Status == "" {
		task.Status = TaskStatusCreated
	}

	return nil
}
