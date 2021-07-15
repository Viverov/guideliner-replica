package token_provider

import (
	"fmt"
)

type NotTokenError struct {
	token string
}

func (e *NotTokenError) Error() string {
	return fmt.Sprintf("Provided string is not a token: %s", e.token)
}

type ExpiredTokenError struct{}

func (e *ExpiredTokenError) Error() string {
	return "Provided token is expired"
}

type UnexpectedTokenError struct {
	token string
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("Can't handle provided token: %s", e.token)
}

type UnexpectedGenerateError struct{}

func (e *UnexpectedGenerateError) Error() string {
	return "Can't generate new token by unexpected error"
}
