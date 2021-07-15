package service

import (
	"errors"
	userEntity "github.com/Viverov/guideliner/internal/domains/user/entity"
	userRepository "github.com/Viverov/guideliner/internal/domains/user/repository"
	tokens "github.com/Viverov/guideliner/internal/domains/user/token_provider"
	"time"
)

type userServiceImpl struct {
	tokenService   tokens.TokenProvider
	userRepository userRepository.UserRepository
	tokenTTL       time.Duration
}

func (u *userServiceImpl) FindById(id uint) (userEntity.UserDTO, error) {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		ID: id,
	})
	if err != nil {
		return nil, processRepositoryError(err)
	}

	return userEntity.NewUserDTO(user.ID(), user.Email()), nil
}

func (u *userServiceImpl) FindByEmail(email string) (userEntity.UserDTO, error) {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		Email: email,
	})
	if err != nil {
		return nil, processRepositoryError(err)
	}

	return userEntity.NewUserDTO(user.ID(), user.Email()), nil
}

func (u *userServiceImpl) Register(email string, password string) (userEntity.UserDTO, error) {
	alreadyExistUser, err := u.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if alreadyExistUser != nil {
		return nil, &EmailAlreadyExistError{}
	}

	user, err := userEntity.CreateUser(email, password)
	if err != nil {
		return nil, &UnexpectedServiceError{}
	}

	id, err := u.userRepository.Insert(user)
	if err != nil {
		return nil, processRepositoryError(err)
	}
	return userEntity.NewUserDTO(id, email), nil
}

func (u *userServiceImpl) ValidateCredentials(email string, password string) (bool, error) {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		Email: email,
	})

	if err != nil {
		return false, processRepositoryError(err)
	}

	isValid := user.ValidatePassword(password)
	return isValid, nil
}

func (u *userServiceImpl) ChangePassword(id uint, newPassword string) error {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		ID: id,
	})
	if err != nil {
		return processRepositoryError(err)
	}

	err = user.SetPassword(newPassword)
	if err != nil {
		return &UnexpectedServiceError{}
	}

	return nil
}

func (u *userServiceImpl) GetToken(userId uint) (string, error) {
	_, err := u.userRepository.FindOne(userRepository.FindCondition{
		ID: userId,
	})
	if err != nil {
		return "", processRepositoryError(err)
	}

	token, err := u.tokenService.GenerateToken(userId, u.tokenTTL)
	if err != nil {
		return "", &UnexpectedServiceError{}
	}

	return token, nil
}

func (u *userServiceImpl) GetUserFromToken(token string) (userEntity.UserDTO, error) {
	claims, err := u.tokenService.ValidateToken(token)
	if err != nil {
		return nil, processTokenError(err)
	}

	userID := claims.UserId
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		ID: userID,
	})
	if err != nil {
		return nil, processRepositoryError(err)
	}

	return userEntity.NewUserDTO(user.ID(), user.Email()), nil
}

func NewUserService(tokenService tokens.TokenProvider, userRepository userRepository.UserRepository, tokenTTL time.Duration) UserService {
	return &userServiceImpl{
		tokenService:   tokenService,
		userRepository: userRepository,
		tokenTTL:       tokenTTL,
	}
}

func processRepositoryError(err error) error {
	var cre *userRepository.CommonRepositoryError
	if errors.As(err, &cre) {
		return &StorageError{storageErrorText: cre.Error()}
	}
	return &UnexpectedServiceError{}
}

func processTokenError(err error) error {
	_, isExpired := err.(*tokens.UnexpectedTokenError)
	_, isNotToken := err.(*tokens.NotTokenError)

	if isExpired || isNotToken {
		return err
	}
	return &UnexpectedServiceError{}
}
