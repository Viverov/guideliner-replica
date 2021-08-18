package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"not null"`
}

func (u User) TableName() string {
	return "users"
}

func createUserTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1621165565",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	}
}
