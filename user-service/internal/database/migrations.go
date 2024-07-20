package database

import (
	"time"

	"github.com/sejamuchhal/task-management/common/utils"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) *gorm.DB {
	logger := utils.Logger
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
		ID: "202407192100",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				ID        uint      `json:"id" gorm:"primary_key"`
				Name      string    `gorm:"size:255" json:"name"`
				Email     string    `gorm:"size:100;not null;unique" json:"email"`
				Password  string    `gorm:"size:100;not null;" json:"password"`
				CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
				UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
			}
			if err := tx.Migrator().CreateTable(&User{}); err != nil {
				tx.Rollback()
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	},
}
