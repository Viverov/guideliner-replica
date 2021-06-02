package user_entity

import "fmt"

const (
	argNameEmail    = "Email"
	argNamePassword = "Password"
)

type User interface {
	// SetEmail sets user email
	SetEmail(email string) error
	// SetPassword hashes & sets new user password
	SetPassword(password string) error
	// ValidatePassword receives raw password as input, returns "true" for correct user password
	ValidatePassword(password string) (isValid bool)
	// ID returns user's ID
	ID() uint
	// Password returns user's password
	Password() string
	// Email returns user's email
	Email() string
}

type EmptyArgError struct {
	argName string
}

func (e *EmptyArgError) Error() string {
	return fmt.Sprintf("%s must not be empty", e.argName)
}
