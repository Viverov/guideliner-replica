package service

import (
	"fmt"
)

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
