package repository

type InvalidFindConditionError struct{}

func NewInvalidFindConditionError() *InvalidFindConditionError {
	return &InvalidFindConditionError{}
}

func (e *InvalidFindConditionError) Error() string {
	return "Error occurred while find user: at least one condition must be defined"
}

type UserAlreadyExistsError struct{}

func NewUserAlreadyExistsError() *UserAlreadyExistsError {
	return &UserAlreadyExistsError{}
}

func (u *UserAlreadyExistsError) Error() string {
	return "The user already exists"
}

type InvalidIdError struct{}

func NewInvalidIdError() *InvalidIdError {
	return &InvalidIdError{}
}

func (i *InvalidIdError) Error() string {
	return "Invalid ID"
}
