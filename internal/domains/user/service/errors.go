package service

type EmailAlreadyExistError struct{}

func NewEmailAlreadyExistError() *EmailAlreadyExistError {
	return &EmailAlreadyExistError{}
}

func (e *EmailAlreadyExistError) Error() string {
	return "Email already exists"
}
