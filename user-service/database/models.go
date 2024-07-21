package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `sql:"size:255;not null";gorm:"primary_key"`
	Name      string    `gorm:"size:255" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Tasks     []Task    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"tasks"`
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return nil
}

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
