package user_dto

type dtoImpl struct {
	id    uint
	email string
}

func NewDTO(id uint, email string) DTO {
	return &dtoImpl{id: id, email: email}
}

func (D *dtoImpl) ID() uint {
	return D.id
}

func (D *dtoImpl) Email() string {
	return D.email
}
