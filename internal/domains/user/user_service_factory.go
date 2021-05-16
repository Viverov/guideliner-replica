package user

import (
	ts "github.com/Viverov/guideliner/internal/domains/user/token_service"
	ur "github.com/Viverov/guideliner/internal/domains/user/user_repository"
	us "github.com/Viverov/guideliner/internal/domains/user/user_service"
	"gorm.io/gorm"
	"time"
)

const defaultTokenDuration = time.Hour * 48

func NewUserService(tokenSecretKey string, db *gorm.DB) us.UserServicer {
	return us.NewUserService(
		ts.NewTokenServiceJWT(tokenSecretKey),
		ur.NewUserRepositoryPostgresql(db),
		defaultTokenDuration,
	)
}
