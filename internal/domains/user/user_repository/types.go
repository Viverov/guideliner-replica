package user_repository

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/domains/user/user_entity"
)

type UserRepository interface {
	FindOne(condition FindCondition) (user_entity.User, error)
	Insert(u user_entity.User) (id uint, err error)
	Update(u user_entity.User) error
}

type FindCondition struct {
	ID    uint
	Email string
}

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
