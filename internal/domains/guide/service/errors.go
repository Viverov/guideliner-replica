package service

import "fmt"

type StorageError struct {
	storageErrorText string
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("storage error: %s", e.storageErrorText)
}

type UnexpectedServiceError struct{}

func (e *UnexpectedServiceError) Error() string {
	return "unexpected error"
}

type GuideNotFoundError struct {
	id uint
}

func (e *GuideNotFoundError) Error() string {
	return fmt.Sprintf("guide with id %d not found", e.id)
}

type InvalidNodesJsonError struct{}

func (e *InvalidNodesJsonError) Error() string {
	return fmt.Sprintf("Can't parse json into correct nodes")
}
