package util

import "fmt"

type NilEntityError struct {
	entityName string
}

func NewNilEntityError(entityName string) *NilEntityError {
	return &NilEntityError{
		entityName: entityName,
	}
}

func (e *NilEntityError) Error() string {
	return fmt.Sprintf("Entity %s can't be nil", e.entityName)
}

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
