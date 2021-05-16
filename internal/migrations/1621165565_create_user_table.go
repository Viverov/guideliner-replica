package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func createUserTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1621165565",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				gorm.Model
				Email    string `gorm:"not null;uniqueIndex"`
				Password string `gorm:"not null"`
			}
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("user")
		},
	}
}
