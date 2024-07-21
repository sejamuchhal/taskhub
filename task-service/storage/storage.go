package storage

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sejamuchhal/taskhub/task-service/common"
)

type Storage struct {
	db *gorm.DB
}

func getConnectionString() string {
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
}

func New() *Storage {
	logger := common.Logger
	logger.Info("Connecting to database")

	dsn := getConnectionString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	storageInstance := &Storage{
		db: db,
	}

	MigrateDB(db)
	return storageInstance
}

func (s *Storage) CreateTask(task *Task) error {
	err := s.db.Create(task).Error
	return err
}

func (s *Storage) GetTaskByID(id string) (*Task, error) {
	var result Task
	err := s.db.Model(&Task{}).First(&result, "id= ?", id).Error
	return &result, err
}

func (s *Storage) ListTasksWithCount(user_id string, limit, offset int) ([]*Task, int64, error) {
	tasks := make([]*Task, 0, limit)
	var count int64
	db := s.db
	if user_id != "" {
		db = db.Where(&Task{UserID: user_id})
	}
	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&tasks).Error
	if err != nil {
		return tasks, count, err
	}
	err = db.Limit(-1).Offset(-1).Count(&count).Error
	return tasks, count, err
}

func (s *Storage) DeleteTask(taskID string) error {
	err := s.db.Delete(&Task{}, "id = ?", taskID).Error
	return err
}

func (s *Storage) UpdateTask(task *Task) error {
	err := s.db.Model(&Task{}).Where("id = ?", task.ID).Updates(task).Error
	return err
}
