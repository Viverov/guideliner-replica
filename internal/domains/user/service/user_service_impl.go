package service

import (
	userEntity "github.com/Viverov/guideliner/internal/domains/user/entity"
	userRepository "github.com/Viverov/guideliner/internal/domains/user/repository"
	tokens "github.com/Viverov/guideliner/internal/domains/user/token_provider"
	"github.com/Viverov/guideliner/internal/domains/util/uservice"
	"strings"
	"time"
)

type userServiceImpl struct {
	tokenProvider  tokens.TokenProvider
	userRepository userRepository.UserRepository
	tokenTTL       time.Duration
}

func NewUserService(tokenService tokens.TokenProvider, userRepository userRepository.UserRepository, tokenTTL time.Duration) UserService {
	return &userServiceImpl{
		tokenProvider:  tokenService,
		userRepository: userRepository,
		tokenTTL:       tokenTTL,
	}
}

func (u *userServiceImpl) FindById(id uint) (userEntity.UserDTO, error) {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		ID: id,
	})
	if err != nil {
		return nil, processRepositoryError(err)
	}
	if user == nil {
		return nil, nil
	}

	return userEntity.NewUserDTOFromEntity(user), nil
}

func (u *userServiceImpl) FindByEmail(email string) (userEntity.UserDTO, error) {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		Email: strings.ToLower(email),
	})
	if err != nil {
		return nil, processRepositoryError(err)
	}
	if user == nil {
		return nil, nil
	}

	return userEntity.NewUserDTOFromEntity(user), nil
}

func (u *userServiceImpl) Register(email string, password string) (userEntity.UserDTO, error) {
	alreadyExistUser, err := u.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if alreadyExistUser != nil {
		return nil, &EmailAlreadyExistError{}
	}

	user, err := userEntity.NewUserWithRawPassword(0, email, password)
	if err != nil {
		return nil, uservice.NewUnexpectedServiceError()
	}

	id, err := u.userRepository.Insert(user)
	if err != nil {
		return nil, processRepositoryError(err)
	}
	err = user.SetID(id)
	if err != nil {
		return nil, err
	}

	return userEntity.NewUserDTOFromEntity(user), nil
}

func (u *userServiceImpl) ValidateCredentials(email string, password string) (bool, error) {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		Email: email,
	})
	if err != nil {
		return false, processRepositoryError(err)
	}
	if user == nil {
		return false, uservice.NewNotFoundError("User", 0)
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
	if user == nil {
		return uservice.NewNotFoundError("User", id)
	}

	err = user.CryptAndSetPassword(newPassword)
	if err != nil {
		return uservice.NewUnexpectedServiceError()
	}

	err = u.userRepository.Update(user)
	if err != nil {
		return processRepositoryError(err)
	}

	return nil
}

func (u *userServiceImpl) GetToken(userId uint) (string, error) {
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		ID: userId,
	})
	if err != nil {
		return "", processRepositoryError(err)
	}
	if user == nil {
		return "", uservice.NewNotFoundError("User", userId)
	}

	token, err := u.tokenProvider.GenerateToken(userId, u.tokenTTL)
	if err != nil {
		return "", uservice.NewUnexpectedServiceError()
	}

	return token, nil
}

func (u *userServiceImpl) GetUserByToken(token string) (userEntity.UserDTO, error) {
	claims, err := u.tokenProvider.ValidateToken(token)
	if err != nil {
		return nil, processTokenError(err)
	}

	userID := claims.UserID
	user, err := u.userRepository.FindOne(userRepository.FindCondition{
		ID: userID,
	})
	if err != nil {
		return nil, processRepositoryError(err)
	}

	return userEntity.NewUserDTOFromEntity(user), nil
}

func processRepositoryError(err error) error {
	switch t := err.(type) {
	case *userRepository.InvalidIdError:
		return uservice.NewStorageError(t.Error())
	default:
		return uservice.CheckDefaultRepoErrors(err)
	}
}

func processTokenError(err error) error {
	switch e := err.(type) {
	case *tokens.UnexpectedTokenError, *tokens.NotTokenError:
		return e
	default:
		return uservice.NewUnexpectedServiceError()
	}
}
