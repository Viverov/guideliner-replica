package user_dto

type DTO interface {
	ID() uint
	Email() string
}
