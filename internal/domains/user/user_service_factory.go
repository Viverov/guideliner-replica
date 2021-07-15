package user

import (
	ur "github.com/Viverov/guideliner/internal/domains/user/repository"
	us "github.com/Viverov/guideliner/internal/domains/user/service"
	ts "github.com/Viverov/guideliner/internal/domains/user/token_provider"
	"gorm.io/gorm"
	"time"
)

const defaultTokenDuration = time.Hour * 48

func BuildUserService(tokenSecretKey string, db *gorm.DB) us.UserService {
	return us.NewUserService(
		ts.NewTokenServiceJWT(tokenSecretKey),
		ur.NewUserRepositoryPostgresql(db),
		defaultTokenDuration,
	)
}
