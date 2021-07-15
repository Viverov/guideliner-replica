package entity

type UserDTO interface {
	ID() uint
	Email() string
}
