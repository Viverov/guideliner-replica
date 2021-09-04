package uservice

import "fmt"

// NotFoundError must be returned on actions (ex: update, patch, delete) with non-existent entities
type NotFoundError struct {
	entityName string
	id         uint
}

func NewNotFoundError(entityName string, id uint) *NotFoundError {
	return &NotFoundError{
		entityName: entityName,
		id:         id,
	}
}

func (e *NotFoundError) Error() string {
	errorText := ""

	if e.entityName != "" {
		errorText += e.entityName
	} else {
		errorText += "Entity"
	}

	if e.id != 0 {
		errorText += fmt.Sprintf(" with ID %d", e.id)
	}

	errorText += " not found"
	return errorText
}
