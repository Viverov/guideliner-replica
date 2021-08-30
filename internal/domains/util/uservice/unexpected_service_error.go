package uservice

// UnexpectedServiceError must be returned on unpredictable actions
type UnexpectedServiceError struct{}

func NewUnexpectedServiceError() *UnexpectedServiceError {
	return &UnexpectedServiceError{}
}

func (e *UnexpectedServiceError) Error() string {
	return "unexpected error"
}
