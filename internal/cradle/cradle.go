package cradle

import (
	"github.com/Viverov/guideliner/internal/config"
	userService "github.com/Viverov/guideliner/internal/domains/user/user_service"
	"gorm.io/gorm"
)

type Cradle struct {
	// Config section
	config *config.Config

	// DB
	sqlDB *gorm.DB

	// Services section
	userService userService.UserService
}

func (c *Cradle) GetConfig() *config.Config {
	return c.config
}

func (c *Cradle) GetSqlDB() *gorm.DB {
	return c.sqlDB
}

func (c *Cradle) GetUserService() userService.UserService {
	return c.userService
}
