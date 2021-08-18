package repository

import "github.com/Viverov/guideliner/internal/domains/user/entity"

type UserRepository interface {
	FindOne(condition FindCondition) (entity.User, error)
	Insert(u entity.User) (id uint, err error)
	Update(u entity.User) error
}

type FindCondition struct {
	ID    uint
	Email string
}
