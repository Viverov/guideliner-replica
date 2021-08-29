package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type GuideModel struct {
	gorm.Model
	Description string `gorm:"not null"`
	NodesJson   string `gorm:"not null"`
}

func (u GuideModel) TableName() string {
	return "guides"
}

func createGuidesTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1629381708",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&GuideModel{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("guides")
		},
	}
}
