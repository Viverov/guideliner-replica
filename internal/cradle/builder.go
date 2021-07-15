package cradle

import (
	"github.com/Viverov/guideliner/internal/config"
	us "github.com/Viverov/guideliner/internal/domains/user/service"
	"gorm.io/gorm"
)

type Builder struct {
	// Config section
	config *config.Config

	// DB
	sqlDB *gorm.DB

	// Services section
	userService us.UserService
}

func (c *Builder) SetConfig(cfg *config.Config) {
	c.config = cfg
}

func (c *Builder) SetSqlDB(db *gorm.DB) {
	c.sqlDB = db
}

func (c *Builder) SetUserService(service us.UserService) {
	c.userService = service
}

func (c *Builder) Build() *Cradle {
	return &Cradle{
		config:      c.config,
		sqlDB:       c.sqlDB,
		userService: c.userService,
	}
}
