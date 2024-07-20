package database

import "time"

type Task struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:50;not null;default:'pending'" json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UserID      uint      `json:"user_id"`
}
