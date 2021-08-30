package repository

type InvalidFindConditionError struct{}

func (e *InvalidFindConditionError) Error() string {
	return "Error occurred while find user: at least one condition must be defined"
}

type UserAlreadyExistsError struct{}

func (u *UserAlreadyExistsError) Error() string {
	return "The user already exists"
}

type InvalidIdError struct{}

func (i *InvalidIdError) Error() string {
	return "Invalid ID"
}
