package urepo

import "fmt"

// EntityNotFoundError must be returned on actions (ex: findOne, update, patch, delete) with non-existent entities
type EntityNotFoundError struct {
	entityName string
	id         uint
}

func NewEntityNotFoundError(entityName string, id uint) *EntityNotFoundError {
	return &EntityNotFoundError{
		entityName: entityName,
		id:         id,
	}
}

func (e *EntityNotFoundError) Error() string {
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

func (e *EntityNotFoundError) EntityName() string {
	return e.entityName
}

func (e *EntityNotFoundError) Id() uint {
	return e.id
}
