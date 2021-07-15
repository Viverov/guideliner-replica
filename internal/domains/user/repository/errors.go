package repository

import (
	"fmt"
)

type InvalidFindConditionError struct{}

func (e *InvalidFindConditionError) Error() string {
	return "Error occurred while find user: at least one condition must be defined"
}

type CommonRepositoryError struct {
	action    string
	errorText string
}

func (c *CommonRepositoryError) Error() string {
	return fmt.Sprintf("Error occured while %s: %s", c.action, c.errorText)
}
