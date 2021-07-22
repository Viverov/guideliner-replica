package entity

import "fmt"

type EmptyArgError struct {
	argName string
}

func (e *EmptyArgError) Error() string {
	return fmt.Sprintf("%s must not be empty", e.argName)
}

type InvalidIdError struct {}

func (e *InvalidIdError) Error() string {
	return fmt.Sprint("ID must be positive")
}
