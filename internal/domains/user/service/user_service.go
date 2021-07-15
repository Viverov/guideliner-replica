package service

import "github.com/Viverov/guideliner/internal/domains/user/entity"

type UserService interface {
	FindById(id uint) (entity.UserDTO, error)
	FindByEmail(email string) (entity.UserDTO, error)
	Register(email string, password string) (entity.UserDTO, error)
	ValidateCredentials(email string, password string) (bool, error)
	ChangePassword(userId uint, newPassword string) error
	GetToken(userId uint) (string, error)
	GetUserFromToken(token string) (entity.UserDTO, error)
}
