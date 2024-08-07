package storage

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
	Active    bool      `gorm:"default:true"`
	Role      string    `gorm:"size:100;default:'user'"`
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

type Session struct {
	ID           string    `gorm:"type:uuid;primary_key;"`
	Email        string    `gorm:"size:100;" json:"email"`
	RefreshToken string    `gorm:"not null;" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	IsBlocked    bool      `gorm:"default:false" json:"is_blocked"`
	BlockedAt    time.Time `json:"blocked_at"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (session *Session) BeforeSave(tx *gorm.DB) (err error) {
	if session.ID == "" {
		session.ID = uuid.NewString()
	}
	return nil
}
