package repository

import "fmt"

type CommonRepositoryError struct {
	action    string
	errorText string
}

func (c *CommonRepositoryError) Error() string {
	return fmt.Sprintf("Error occured while %s: %s", c.action, c.errorText)
}

type InvalidFindConditionError struct{}

func (e *InvalidFindConditionError) Error() string {
	return "Error occurred while find guides: search must be defined"
}

type GuideNotFoundError struct{}

func (e *GuideNotFoundError) Error() string {
	return "Guide with provided ID not found"
}
