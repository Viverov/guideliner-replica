package user_repository

import (
	userEntity "github.com/Viverov/guideliner/internal/domains/user/user_entity"
	"gorm.io/gorm"
)

type userModel struct {
	gorm.Model
	Email    string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"not null"`
}

type userRepositoryPostgresql struct {
	db *gorm.DB
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

	result := r.db.Find(um)
	if result.Error != nil {
		return nil, &CommonRepositoryError{
			action:    "find",
			errorText: result.Error.Error(),
		}
	}

	user, err := userEntity.NewUser(um.ID, um.Email, um.Password)
	return user, err
}

func (r *userRepositoryPostgresql) Insert(u userEntity.User) (id uint, err error) {
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

func NewUserRepositoryPostgresql(db *gorm.DB) UserRepositorer {
	return &userRepositoryPostgresql{
		db: db,
	}
}
