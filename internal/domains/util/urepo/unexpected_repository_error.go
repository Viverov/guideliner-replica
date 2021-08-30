package urepo

import "fmt"

// UnexpectedRepositoryError can be returned on any unexpected actions with repo, like broken DB connection
type UnexpectedRepositoryError struct {
	action    string
	errorText string
}

func NewUnexpectedRepositoryError(action string, errorText string) *UnexpectedRepositoryError {
	return &UnexpectedRepositoryError{
		action:    action,
		errorText: errorText,
	}
}

func (c *UnexpectedRepositoryError) Error() string {
	return fmt.Sprintf("Error occured while %s: %s", c.action, c.errorText)
}
