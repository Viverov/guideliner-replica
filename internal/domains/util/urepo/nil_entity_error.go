package urepo

import "fmt"

// NilEntityError must be returned on actions with nil entity (ex: insert, update)
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
