package repository

import (
	"fmt"
)

type InvalidFindConditionError struct{}

func (e *InvalidFindConditionError) Error() string {
	return "Error occurred while find user: at least one condition must be defined"
}

type CommonRepositoryError struct {
	Action    string
	ErrorText string
}

func (c *CommonRepositoryError) Error() string {
	return fmt.Sprintf("Error occured while %s: %s", c.Action, c.ErrorText)
}

type UserAlreadyExistsError struct{}

func (u *UserAlreadyExistsError) Error() string {
	return "The user already exists"
}

type UserNotFoundError struct{}

func (u *UserNotFoundError) Error() string {
	return "User not found"
}

type InvalidIdError struct{}

func (i *InvalidIdError) Error() string {
	return "Invalid ID"
}

func (i *InvalidIdError) RuntimeError() {
	panic("implement me")
}
