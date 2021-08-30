package service

type EmailAlreadyExistError struct{}

func (e *EmailAlreadyExistError) Error() string {
	return "Email already exists"
}
