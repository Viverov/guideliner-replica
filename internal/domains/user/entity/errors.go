package entity

import "fmt"

type EmptyArgError struct {
	argName string
}

func (e *EmptyArgError) Error() string {
	return fmt.Sprintf("%s must not be empty", e.argName)
}
