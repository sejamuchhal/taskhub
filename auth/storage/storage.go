package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sejamuchhal/taskhub/auth/common"
)

type StorageInterface interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	GetSessionByID(sessionID string) (*Session, error)
	UpdateSession(session *Session) error
}

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

func (s *Storage) CreateUser(user *User) error {
	err := s.db.Create(user).Error
	return err
}

func (s *Storage) GetUserByEmail(email string) (*User, error) {
	var result User
	err := s.db.Model(&User{}).First(&result, "email= ?", email).Error
	return &result, err
}

func (s *Storage) CreateSession(session *Session) error {
	err := s.db.Create(session).Error

	return err
}

func (s *Storage) GetSessionByID(id string) (*Session, error) {
	var result Session
	err := s.db.Model(&Session{}).First(&result, "id= ?", id).Error
	return &result, err
}

func (s *Storage) BlockSessionByID(id string) error {
	err := s.db.Model(&Session{}).Where("id= ?", id).Updates(map[string]interface{}{
		"is_blocked": true,
		"blocked_at": time.Now(),
	}).Error
	return err
}

func (s *Storage) DeleteSessionByID(id string) error {
	err := s.db.Model(&Session{}).Delete("id= ?", id).Error
	return err
}
