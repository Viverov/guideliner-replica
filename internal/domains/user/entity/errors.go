package entity

import "fmt"

type EmptyArgError struct {
	argName string
}

func NewEmptyArgError(argName string) *EmptyArgError {
	return &EmptyArgError{argName: argName}
}

func (e *EmptyArgError) Error() string {
	return fmt.Sprintf("%s must not be empty", e.argName)
}

type InvalidIdError struct{}

func NewInvalidIdError() *InvalidIdError {
	return &InvalidIdError{}
}

func (e *InvalidIdError) Error() string {
	return "ID must be positive"
}
