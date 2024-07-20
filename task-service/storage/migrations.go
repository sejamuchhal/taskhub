package storage

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sejamuchhal/taskhub/task-service/common"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) *gorm.DB {
	logger := common.Logger
	logger = logger.WithField("method", "MigrateDB")
	dbObj, err := db.DB()
	if err != nil {
		logger.WithError(err).Error("Failed to get db")
	}
	if err := dbObj.Ping(); err != nil {
		logger.WithError(err).Fatalln("Could not ping DB")
	}
	logger.Info("running migrations")
	options := &gormigrate.Options{
		IDColumnName:   "id",
		UseTransaction: true,
	}

	migration := gormigrate.New(db, options, migrations)
	if err := migration.Migrate(); err != nil {
		logger.WithError(err).Fatalln("failed to migrate database")
	}

	logger.Info("migrations completed")
	return db
}

var migrations = []*gormigrate.Migration{
	{
		ID: "202407200130",
		Migrate: func(tx *gorm.DB) error {
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

			if err := tx.Migrator().CreateTable(&Task{}); err != nil {
				tx.Rollback()
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("tasks")
		},
	},
}
