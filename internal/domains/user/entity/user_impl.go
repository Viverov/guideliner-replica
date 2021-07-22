package entity

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type userImpl struct {
	id       uint
	email    string
	password string
}

func (u *userImpl) SetID(id uint) error {
	if id == 0 {
		return &InvalidIdError{}
	}
	u.id = id

	return nil
}

func (u *userImpl) SetPassword(password string) error {
	if len(password) == 0 {
		return &EmptyArgError{argName: argNamePassword}
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("something gone wrong on password generate: %s", err.Error())
	}

	u.password = string(bytes)
	return nil
}

func (u *userImpl) ValidatePassword(password string) (isValid bool) {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}

func (u *userImpl) SetEmail(email string) error {
	if len(email) == 0 {
		return &EmptyArgError{argName: argNameEmail}
	}
	u.email = strings.ToLower(email)
	return nil
}

func (u *userImpl) ID() uint {
	return u.id
}

func (u *userImpl) Password() string {
	return u.password
}

func (u *userImpl) Email() string {
	return u.email
}

// CreateUser receives email and unhashed password as input values, returns user object without ID (will be set later in db / repository)
func CreateUser(email string, password string) (User, error) {
	if len(email) == 0 {
		return nil, &EmptyArgError{argName: argNameEmail}
	}
	if len(password) == 0 {
		return nil, &EmptyArgError{argName: argNamePassword}
	}

	u := &userImpl{}
	err := u.SetEmail(email)
	if err != nil {
		return nil, err
	}
	err = u.SetPassword(password)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// NewUser is user constructor. Use ID == 0 for new user.
func NewUser(id uint, email string, password string) (User, error) {
	if len(email) == 0 {
		return nil, &EmptyArgError{argName: argNameEmail}
	}
	if len(password) == 0 {
		return nil, &EmptyArgError{argName: argNamePassword}
	}
	return &userImpl{id, strings.ToLower(email), password}, nil
}
