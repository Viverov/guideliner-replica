package entity

type User interface {
	// SetID sets user's ID
	SetID(id uint) error
	// SetEmail sets user's email
	SetEmail(email string) error
	// SetPassword hashes & sets new user's password
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
