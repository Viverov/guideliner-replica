package entity

type userDTOImpl struct {
	id    uint
	email string
}

func NewUserDTO(id uint, email string) UserDTO {
	return &userDTOImpl{id: id, email: email}
}

func (D *userDTOImpl) ID() uint {
	return D.id
}

func (D *userDTOImpl) Email() string {
	return D.email
}
