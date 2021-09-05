package cradle

import (
	"github.com/Viverov/guideliner/internal/config"
	gs "github.com/Viverov/guideliner/internal/domains/guide/service"
	us "github.com/Viverov/guideliner/internal/domains/user/service"
	"gorm.io/gorm"
)

type Cradle struct {
	// Config section
	config *config.Config

	// DB
	sqlDB *gorm.DB

	// Services section
	userService  us.UserService
	guideService gs.GuideService
}

func (c *Cradle) GetConfig() *config.Config {
	return c.config
}

func (c *Cradle) GetSqlDB() *gorm.DB {
	return c.sqlDB
}

func (c *Cradle) GetUserService() us.UserService {
	return c.userService
}

func (c *Cradle) GetGuideService() gs.GuideService {
	return c.guideService
}
