package repository

import (
	"errors"
	userEntity "github.com/Viverov/guideliner/internal/domains/user/entity"
	urepo "github.com/Viverov/guideliner/internal/domains/util/urepo"
	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	Email    string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"not null"`
}

func (u userModel) TableName() string {
	return "users"
}

type userRepositoryPostgresql struct {
	db *gorm.DB
}

func NewUserRepositoryPostgresql(db *gorm.DB) *userRepositoryPostgresql {
	return &userRepositoryPostgresql{
		db: db,
	}
}

func (r *userRepositoryPostgresql) FindOne(condition FindCondition) (userEntity.User, error) {
	if condition.ID == 0 && condition.Email == "" {
		return nil, NewInvalidFindConditionError()
	}

	um := &userModel{
		Model: gorm.Model{
			ID: condition.ID,
		},
		Email: condition.Email,
	}

	result := r.db.Where(um).First(um)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, urepo.NewUnexpectedRepositoryError("Find", result.Error.Error())
	}

	user, err := userEntity.NewUser(um.ID, um.Email, um.Password)
	return user, err
}

func (r *userRepositoryPostgresql) Insert(u userEntity.User) (id uint, err error) {
	alreadyExistsUser, err := r.FindOne(FindCondition{Email: u.Email()})
	if err != nil {
		return 0, err
	}
	if alreadyExistsUser != nil {
		return 0, NewUserAlreadyExistsError()
	}

	um := &userModel{
		Email:    u.Email(),
		Password: u.Password(),
	}
	result := r.db.Create(um)

	if result.Error != nil {
		return 0, urepo.NewUnexpectedRepositoryError("Create", result.Error.Error())
	}

	return um.ID, nil
}

func (r *userRepositoryPostgresql) Update(u userEntity.User) error {
	if u.ID() == 0 {
		return NewInvalidIdError()
	}

	user, err := r.FindOne(FindCondition{ID: u.ID()})
	if err != nil {
		return err
	}
	if user == nil {
		return urepo.NewEntityNotFoundError("User", u.ID())
	}

	um := &userModel{
		Model: gorm.Model{
			ID: u.ID(),
		},
		Email:    u.Email(),
		Password: u.Password(),
	}

	result := r.db.Save(um)
	if result.Error != nil {
		return urepo.NewUnexpectedRepositoryError("Update", result.Error.Error())
	}

	return nil
}
