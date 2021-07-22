package entity

type userDTOImpl struct {
	id    uint
	email string
}

func NewUserDTOFromEntity(user User) *userDTOImpl {
	return &userDTOImpl{
		id:    user.ID(),
		email: user.Email(),
	}
}

func (D *userDTOImpl) ID() uint {
	return D.id
}

func (D *userDTOImpl) Email() string {
	return D.email
}
