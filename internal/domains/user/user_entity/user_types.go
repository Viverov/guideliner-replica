package user_entity

import "fmt"

const (
	argNameEmail    = "Email"
	argNamePassword = "Password"
)

type User interface {
	// Set user email
	SetEmail(email string) error
	// Hash & set new user password
	SetPassword(password string) error
	// Validate user password. Get unhashed password as input, return "true" for correct user password
	ValidatePassword(password string) (isValid bool)
	// Get user ID
	ID() uint
	// Get user password
	Password() string
	// Get user email
	Email() string
}

type EmptyArgError struct {
	argName string
}

func (e *EmptyArgError) Error() string {
	return fmt.Sprintf("%s must not be empty", e.argName)
}
