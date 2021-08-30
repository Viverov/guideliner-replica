package entity

import "fmt"

type UnexpectedGuideError struct {
	info string
}

func NewUnexpectedGuideError(info string) *UnexpectedGuideError {
	return &UnexpectedGuideError{info: info}
}

func (u *UnexpectedGuideError) Error() string {
	return fmt.Sprintf("unexpected guide error: %s", u.info)
}

type InvalidJsonError struct{}

func NewInvalidJsonError() *InvalidJsonError {
	return &InvalidJsonError{}
}

func (e *InvalidJsonError) Error() string {
	return "invalid json passed"
}

type InvalidIdError struct{}

func NewInvalidIdError() *InvalidIdError {
	return &InvalidIdError{}
}

func (e *InvalidIdError) Error() string {
	return "id must be above zero"
}
