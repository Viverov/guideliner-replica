package uservice

import "fmt"

// StorageError must be returned on urepo.UnexpectedRepositoryError
type StorageError struct {
	storageErrorText string
}

func NewStorageError(storageErrorText string) *StorageError {
	return &StorageError{storageErrorText: storageErrorText}
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("storage error: %s", e.storageErrorText)
}
