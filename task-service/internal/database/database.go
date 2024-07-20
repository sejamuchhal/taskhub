package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sejamuchhal/task-management/common"
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
