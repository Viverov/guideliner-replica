package user_service

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/domains/user/user_entity"
)

type UserService interface {
	FindById(id uint) (user_entity.DTO, error)
	FindByEmail(email string) (user_entity.DTO, error)
	Register(email string, password string) (user_entity.DTO, error)
	ValidateCredentials(email string, password string) (bool, error)
	ChangePassword(userId uint, newPassword string) error
	GetToken(userId uint) (string, error)
	GetUserFromToken(token string) (user_entity.DTO, error)
}

type InvalidNewArgsError struct {
	argsName []string
}

func (e *InvalidNewArgsError) Error() string {
	return fmt.Sprintf("Invalid new args: %v", e.argsName)
}

type EmailAlreadyExistError struct{}

func (e *EmailAlreadyExistError) Error() string {
	return "Email already exists"
}

type UnexpectedServiceError struct{}

func (i *UnexpectedServiceError) Error() string {
	return "Unexpected error"
}

type StorageError struct {
	storageErrorText string
}

func (s *StorageError) Error() string {
	return fmt.Sprintf("Storage error: %s", s.storageErrorText)
}
