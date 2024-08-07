package storage

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sejamuchhal/taskhub/auth/common"
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
		ID: "202407211130",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				ID        string    `gorm:"size:255;not null;primary_key" json:"id"`
				Name      string    `gorm:"size:255" json:"name"`
				Email     string    `gorm:"size:100;not null;unique" json:"email"`
				Password  string    `gorm:"size:100;not null" json:"password"`
				CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
				UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
				Tasks     []Task    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"tasks"`
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
	{
		ID: "202407231130",
		Migrate: func(tx *gorm.DB) error {
			log.Println("add active field to user Model")
			type User struct {
				Active bool `gorm:"default:true"`
			}
			if err := tx.Migrator().AddColumn(&User{}, "Active"); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			type User struct {
				Active bool
			}
			return tx.Migrator().DropColumn(&User{}, "Active")
		},
	},
	{
		ID: "202407231200",
		Migrate: func(tx *gorm.DB) error {
			log.Println("add role field to user Model")
			type User struct {
				Role string `gorm:"size:100;default:'user'"`
			}
			if err := tx.Migrator().AddColumn(&User{}, "Role"); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			type User struct {
				Role string
			}
			return tx.Migrator().DropColumn(&User{}, "Role")
		},
	},
	{
		ID: "202407271230",
		Migrate: func(tx *gorm.DB) error {
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
			if err := tx.Migrator().CreateTable(&Session{}); err != nil {
				tx.Rollback()
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("sessions")
		},
	},
}
