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

// NewUser is user constructor. Use ID == 0 for new user, not saved in DB.
func NewUser(id uint, email string, hashedPassword string) (*userImpl, error) {
	if len(email) == 0 {
		return nil, &EmptyArgError{argName: argNameEmail}
	}
	if len(hashedPassword) == 0 {
		return nil, &EmptyArgError{argName: argNamePassword}
	}
	return &userImpl{id, strings.ToLower(email), hashedPassword}, nil
}

// NewUserWithRawPassword is user constructor for non-crypted password. Use ID == 0 for new user, not saved in DB
func NewUserWithRawPassword(id uint, email string, rawPassword string) (*userImpl, error) {
	u, err := NewUser(id, email, rawPassword)
	if err != nil {
		return nil, err
	}

	err = u.CryptAndSetPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *userImpl) SetID(id uint) error {
	if id == 0 {
		return &InvalidIdError{}
	}
	u.id = id

	return nil
}

func (u *userImpl) CryptAndSetPassword(password string) error {
	cryptedPassword, err := cryptPassword(password)
	if err != nil {
		return err
	}

	u.password = cryptedPassword
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

func cryptPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", &EmptyArgError{argName: argNamePassword}
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("something gone wrong on password generate: %s", err.Error())
	}

	return string(bytes), nil
}
