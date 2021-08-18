package service

import (
	"fmt"
)

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

type UserNotFoundError struct {
	id uint
}

func (e *UserNotFoundError) Error() string {
	if e.id != 0 {
		return fmt.Sprintf("User with id %d not found", e.id)
	}

	return "User not found"
}
