package repository

import (
	"errors"
	userEntity "github.com/Viverov/guideliner/internal/domains/user/entity"
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
		return nil, &InvalidFindConditionError{}
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

		return nil, &CommonRepositoryError{
			action:    "find",
			errorText: result.Error.Error(),
		}
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
		return 0, &UserAlreadyExistsError{}
	}

	um := &userModel{
		Email:    u.Email(),
		Password: u.Password(),
	}
	result := r.db.Create(um)

	if result.Error != nil {
		return 0, &CommonRepositoryError{
			action:    "create",
			errorText: result.Error.Error(),
		}
	}

	return um.ID, nil
}

func (r *userRepositoryPostgresql) Update(u userEntity.User) error {
	um := &userModel{
		Model: gorm.Model{
			ID: u.ID(),
		},
		Email:    u.Email(),
		Password: u.Password(),
	}

	result := r.db.Save(um)

	if result.Error != nil {
		return &CommonRepositoryError{
			action:    "update",
			errorText: result.Error.Error(),
		}
	}

	return nil
}
