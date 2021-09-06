package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type GuideModelWithCreatorId struct {
	gorm.Model
	Description string `gorm:"not null"`
	NodesJson   string `gorm:"not null"`
	CreatorID   uint
}

func (u GuideModelWithCreatorId) TableName() string {
	return "guides"
}

func addCreatorColumnIntoGuides() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1630955515",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AddColumn(&GuideModelWithCreatorId{}, "CreatorID")
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropColumn(&GuideModelWithCreatorId{}, "CreatorID")
		},
	}
}
