// +build unit

package service

import (
	"github.com/Viverov/guideliner/internal/domains/user/entity"
	userRepository "github.com/Viverov/guideliner/internal/domains/user/repository"
	"github.com/Viverov/guideliner/internal/domains/user/service/mocks"
	tokenprovider "github.com/Viverov/guideliner/internal/domains/user/token_provider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		tokenTTL time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should return new UserService",
			args: args{tokenTTL: time.Duration(1000)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urCtrl, userRepositoryMock, _, tokenProviderMock := prepareMocks(t)

			got := NewUserService(tokenProviderMock, userRepositoryMock, tt.args.tokenTTL)
			assert.NotNil(t, got)

			urCtrl.Finish()
		})
	}
}

func Test_userServiceImpl_FindById(t *testing.T) {
	type args struct {
		id uint
	}
	tests := []struct {
		name               string
		args               args
		userFromRepository entity.User
		errFromRepository  error
		want               entity.UserDTO
		wantErr            error
	}{
		{
			name:               "Should return nil for nil return from repository",
			args:               args{10},
			userFromRepository: nil,
			errFromRepository:  nil,
			want:               nil,
			wantErr:            nil,
		},
		{
			name:               "Should return storage error for invalid ID error",
			args:               args{0},
			userFromRepository: nil,
			errFromRepository:  &userRepository.InvalidIdError{},
			want:               nil,
			wantErr:            &StorageError{storageErrorText: (&userRepository.InvalidIdError{}).Error()},
		},
		{
			name:               "Should return storage error for common repository error",
			args:               args{0},
			userFromRepository: nil,
			errFromRepository:  &userRepository.CommonRepositoryError{Action: "test", ErrorText: "errT"},
			want:               nil,
			wantErr:            &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{Action: "test", ErrorText: "errT"}).Error()},
		},
		{
			name:               "Should return user's DTO",
			args:               args{0},
			userFromRepository: func() entity.User { us, _ := entity.NewUser(10, "some@email.com", "pass123!@#"); return us }(),
			errFromRepository:  nil,
			want:               entity.NewUserDTO(10, "some@email.com"),
			wantErr:            nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup (before each)
			urCtrl, userRepositoryMock, _, tokenProviderMock := prepareMocks(t)
			u := &userServiceImpl{
				tokenProvider:  tokenProviderMock,
				userRepository: userRepositoryMock,
				tokenTTL:       time.Duration(10000),
			}

			userRepositoryMock.EXPECT().FindOne(userRepository.FindCondition{
				ID: tt.args.id,
			}).Return(tt.userFromRepository, tt.errFromRepository)

			// Actions
			got, err := u.FindById(tt.args.id)

			// Check
			assertUserDTOAndError(t, tt.want, tt.wantErr, got, err)
			urCtrl.Finish()
		})
	}
}

func Test_userServiceImpl_FindByEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name               string
		args               args
		userFromRepository entity.User
		errFromRepository  error
		want               entity.UserDTO
		wantErr            error
	}{
		{
			name:               "Should return nil for nil return from repository",
			args:               args{"some@email.com"},
			userFromRepository: nil,
			errFromRepository:  nil,
			want:               nil,
			wantErr:            nil,
		},
		{
			name:               "Should return storage error for common repository error",
			args:               args{"some@email.com"},
			userFromRepository: nil,
			errFromRepository:  &userRepository.CommonRepositoryError{Action: "test", ErrorText: "errT"},
			want:               nil,
			wantErr:            &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{Action: "test", ErrorText: "errT"}).Error()},
		},
		{
			name:               "Should lowercase email",
			args:               args{"SoMe@eMaIl.CoM"},
			userFromRepository: nil,
			errFromRepository:  nil,
			want:               nil,
			wantErr:            nil,
		},
		{
			name:               "Should return user's DTO",
			args:               args{"some@email.com"},
			userFromRepository: func() entity.User { us, _ := entity.NewUser(10, "some@email.com", "pass123!@#"); return us }(),
			errFromRepository:  nil,
			want:               entity.NewUserDTO(10, "some@email.com"),
			wantErr:            nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup (before each)
			urCtrl, userRepositoryMock, _, tokenProviderMock := prepareMocks(t)
			u := &userServiceImpl{
				tokenProvider:  tokenProviderMock,
				userRepository: userRepositoryMock,
				tokenTTL:       time.Duration(10000),
			}

			userRepositoryMock.EXPECT().FindOne(userRepository.FindCondition{
				Email: strings.ToLower(tt.args.email),
			}).Return(tt.userFromRepository, tt.errFromRepository)

			// Actions
			got, err := u.FindByEmail(tt.args.email)

			// Check
			assertUserDTOAndError(t, tt.want, tt.wantErr, got, err)
			urCtrl.Finish()
		})
	}
}

func Test_userServiceImpl_Register(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name                    string
		args                    args
		alreadyExistsUser       entity.User
		errFromRepositoryOnFind error
		idFromInsert            uint
		want                    entity.UserDTO
		wantErr                 error
	}{
		{
			name:                    "Should register new user",
			args:                    args{email: "someemail@email.com", password: "abcdefasdasdasd"},
			alreadyExistsUser:       nil,
			errFromRepositoryOnFind: nil,
			idFromInsert:            10,
			want:                    entity.NewUserDTO(10, "someemail@email.com"),
			wantErr:                 nil,
		},
		{
			name:                    "Should return error for already existing user",
			args:                    args{email: "someemail@email.com", password: "abcdefasdasdasd"},
			alreadyExistsUser:       func() entity.User { u, _ := entity.NewUser(10, "someemail@email.com", "abcdef"); return u }(),
			errFromRepositoryOnFind: nil,
			idFromInsert:            0,
			want:                    nil,
			wantErr:                 &EmailAlreadyExistError{},
		},
		{
			name:              "Should return error on repository error",
			args:              args{email: "someemail@email.com", password: "abcdefasdasdasd"},
			alreadyExistsUser: nil,
			errFromRepositoryOnFind: &userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			},
			idFromInsert: 0,
			want:         nil,
			wantErr: &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			}).Error()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup (before each)
			urCtrl, userRepositoryMock, _, tokenProviderMock := prepareMocks(t)
			u := &userServiceImpl{
				tokenProvider:  tokenProviderMock,
				userRepository: userRepositoryMock,
				tokenTTL:       time.Duration(10000),
			}

			userRepositoryMock.
				EXPECT().
				FindOne(userRepository.FindCondition{Email: tt.args.email}).
				Return(tt.alreadyExistsUser, tt.errFromRepositoryOnFind)

			if tt.idFromInsert != 0 {
				userRepositoryMock.
					EXPECT().
					Insert(gomock.Any()).
					Return(tt.idFromInsert, nil)
			}

			// Actions
			got, err := u.Register(tt.args.email, tt.args.password)

			// Check
			assertUserDTOAndError(t, tt.want, tt.wantErr, got, err)
			urCtrl.Finish()
		})
	}
}

func Test_userServiceImpl_ChangePassword(t *testing.T) {
	type args struct {
		id          uint
		newPassword string
	}
	tests := []struct {
		name                      string
		args                      args
		alreadyExistsUser         entity.User
		errFromRepositoryOnFind   error
		errFromRepositoryOnUpdate error
		wantErr                   error
	}{
		{
			name: "Should change password",
			args: args{
				id:          10,
				newPassword: "abcdef",
			},
			alreadyExistsUser:         func() entity.User { u, _ := entity.NewUserWithRawPassword(10, "some@email.com", "qwerty"); return u }(),
			errFromRepositoryOnFind:   nil,
			errFromRepositoryOnUpdate: nil,
			wantErr:                   nil,
		},
		{
			name: "Should return error on undefined user",
			args: args{
				id:          10,
				newPassword: "abcdef",
			},
			alreadyExistsUser:         nil,
			errFromRepositoryOnFind:   nil,
			errFromRepositoryOnUpdate: nil,
			wantErr:                   &UserNotFoundError{id: 10},
		},
		{
			name: "Should return error on find in db related error",
			args: args{
				id:          10,
				newPassword: "abcdef",
			},
			alreadyExistsUser: nil,
			errFromRepositoryOnFind: &userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			},
			errFromRepositoryOnUpdate: nil,
			wantErr: &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			}).Error()},
		},
		{
			name: "Should return error on update in db related error",
			args: args{
				id:          10,
				newPassword: "abcdef",
			},
			alreadyExistsUser:       func() entity.User { u, _ := entity.NewUserWithRawPassword(10, "some@email.com", "qwerty"); return u }(),
			errFromRepositoryOnFind: nil,
			errFromRepositoryOnUpdate: &userRepository.CommonRepositoryError{
				Action:    "update",
				ErrorText: "test",
			},
			wantErr: &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{
				Action:    "update",
				ErrorText: "test",
			}).Error()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup (before each)
			urCtrl, userRepositoryMock, _, tokenProviderMock := prepareMocks(t)

			u := &userServiceImpl{
				tokenProvider:  tokenProviderMock,
				userRepository: userRepositoryMock,
				tokenTTL:       time.Duration(10000),
			}

			var sourcePassword string
			userRepositoryMock.
				EXPECT().
				FindOne(userRepository.FindCondition{ID: tt.args.id}).
				DoAndReturn(func(fc userRepository.FindCondition) (entity.User, error) {
					if tt.alreadyExistsUser != nil {
						sourcePassword = tt.alreadyExistsUser.Password()
					}
					return tt.alreadyExistsUser, tt.errFromRepositoryOnFind
				})

			var capturedPassword string
			if tt.alreadyExistsUser != nil && tt.errFromRepositoryOnFind == nil {
				userRepositoryMock.
					EXPECT().
					Update(gomock.Any()).
					DoAndReturn(func(u entity.User) error {
						capturedPassword = u.Password()
						return tt.errFromRepositoryOnUpdate
					})
			}

			// Actions
			err := u.ChangePassword(tt.args.id, tt.args.newPassword)

			// Check
			assert.Equal(t, tt.wantErr, err)
			// Should set new password to user
			if tt.errFromRepositoryOnFind == nil && tt.errFromRepositoryOnUpdate == nil && tt.alreadyExistsUser != nil {
				assert.NotEqual(t, sourcePassword, capturedPassword)
				assert.Equal(t, tt.alreadyExistsUser.Password(), capturedPassword)
			}

			urCtrl.Finish()
		})
	}
}

func Test_userServiceImpl_GetToken(t *testing.T) {
	type args struct {
		userId uint
	}
	tests := []struct {
		name                    string
		args                    args
		alreadyExistsUser       entity.User
		errFromRepositoryOnFind error
		tokenFromProvider       string
		errFromTokenProvider    error
		want                    string
		wantErr                 error
	}{
		{
			name: "Should get token",
			args: args{
				userId: 10,
			},
			alreadyExistsUser:       func() entity.User { u, _ := entity.NewUserWithRawPassword(10, "some@email.com", "123"); return u }(),
			errFromRepositoryOnFind: nil,
			tokenFromProvider:       "abcdef",
			errFromTokenProvider:    nil,
			want:                    "abcdef",
			wantErr:                 nil,
		},
		{
			name: "Should return error on undefined user",
			args: args{
				userId: 10,
			},
			alreadyExistsUser:       nil,
			errFromRepositoryOnFind: nil,
			tokenFromProvider:       "abcdef",
			errFromTokenProvider:    nil,
			want:                    "",
			wantErr:                 &UserNotFoundError{id: 10},
		},
		{
			name: "Should return error on find in db related error",
			args: args{
				userId: 10,
			},
			alreadyExistsUser: nil,
			errFromRepositoryOnFind: &userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			},
			tokenFromProvider:    "abcdef",
			errFromTokenProvider: nil,
			want:                 "",
			wantErr: &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			}).Error()},
		},
		{
			name: "Should return error on token provider error",
			args: args{
				userId: 10,
			},
			alreadyExistsUser:       func() entity.User { u, _ := entity.NewUserWithRawPassword(10, "some@email.com", "123"); return u }(),
			errFromRepositoryOnFind: nil,
			tokenFromProvider:       "abcdef",
			errFromTokenProvider:    &tokenprovider.UnexpectedGenerateError{},
			want:                    "",
			wantErr:                 &UnexpectedServiceError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup (before each)
			urCtrl, userRepositoryMock, tkCtrl, tokenProviderMock := prepareMocks(t)

			u := &userServiceImpl{
				tokenProvider:  tokenProviderMock,
				userRepository: userRepositoryMock,
				tokenTTL:       time.Duration(10000),
			}

			userRepositoryMock.
				EXPECT().
				FindOne(userRepository.FindCondition{ID: tt.args.userId}).
				Return(tt.alreadyExistsUser, tt.errFromRepositoryOnFind)

			if tt.alreadyExistsUser != nil && tt.errFromRepositoryOnFind == nil {
				tokenProviderMock.
					EXPECT().
					GenerateToken(tt.args.userId, u.tokenTTL).
					Return(tt.tokenFromProvider, tt.errFromTokenProvider)
			}

			got, err := u.GetToken(tt.args.userId)
			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			urCtrl.Finish()
			tkCtrl.Finish()
		})
	}
}

func Test_userServiceImpl_GetUserByToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name                    string
		args                    args
		claims                  *tokenprovider.AuthClaims
		errFromTokenProvider    error
		userFromDB              entity.User
		errFromRepositoryOnFind error
		want                    entity.UserDTO
		wantErr                 error
	}{
		{
			name: "Should find user by token",
			args: args{token: "abcd"},
			claims: &tokenprovider.AuthClaims{
				UserID: 10,
			},
			errFromTokenProvider:    nil,
			userFromDB:              func() entity.User { u, _ := entity.NewUser(10, "some@email.com", "123"); return u }(),
			errFromRepositoryOnFind: nil,
			want:                    entity.NewUserDTO(10, "some@email.com"),
			wantErr:                 nil,
		},
		{
			name:                    "Should return error on token provider error",
			args:                    args{token: "abcd"},
			claims:                  nil,
			errFromTokenProvider:    &tokenprovider.NotTokenError{},
			userFromDB:              nil,
			errFromRepositoryOnFind: nil,
			want:                    nil,
			wantErr:                 &tokenprovider.NotTokenError{},
		},
		{
			name:                 "Should return error on error from repository",
			args:                 args{token: "abcd"},
			claims:               &tokenprovider.AuthClaims{UserID: 10},
			errFromTokenProvider: nil,
			userFromDB:           nil,
			errFromRepositoryOnFind: &userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			},
			want: nil,
			wantErr: &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			}).Error()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup (before each)
			urCtrl, userRepositoryMock, tpCtrl, tokenProviderMock := prepareMocks(t)

			u := &userServiceImpl{
				tokenProvider:  tokenProviderMock,
				userRepository: userRepositoryMock,
				tokenTTL:       time.Duration(10000),
			}

			tokenProviderMock.
				EXPECT().
				ValidateToken(tt.args.token).
				Return(tt.claims, tt.errFromTokenProvider)

			if tt.errFromTokenProvider == nil {
				userRepositoryMock.
					EXPECT().
					FindOne(userRepository.FindCondition{ID: tt.claims.UserID}).
					Return(tt.userFromDB, tt.errFromRepositoryOnFind)
			}

			got, err := u.GetUserByToken(tt.args.token)

			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			urCtrl.Finish()
			tpCtrl.Finish()
		})
	}
}

func Test_userServiceImpl_ValidateCredentials(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name                    string
		args                    args
		userInDB                entity.User
		errFromRepositoryOnFind error
		want                    bool
		wantErr                 error
	}{
		{
			name: "Should return 'true' on valid credentials",
			args: args{
				email:    "some@email.com",
				password: "123123",
			},
			userInDB:                func() entity.User { u, _ := entity.NewUserWithRawPassword(10, "some@email.com", "123123"); return u }(),
			errFromRepositoryOnFind: nil,
			want:                    true,
			wantErr:                 nil,
		},
		{
			name: "Should return 'false' on invalid credentials",
			args: args{
				email:    "some@email.com",
				password: "123123",
			},
			userInDB: func() entity.User {
				u, _ := entity.NewUserWithRawPassword(10, "some@email.com", "invalid password")
				return u
			}(),
			errFromRepositoryOnFind: nil,
			want:                    false,
			wantErr:                 nil,
		},
		{
			name: "Should return error on undefined email",
			args: args{
				email:    "indefined@email.com",
				password: "123123",
			},
			userInDB:                nil,
			errFromRepositoryOnFind: nil,
			want:                    false,
			wantErr:                 &UserNotFoundError{},
		},
		{
			name: "Should return error on error from repository",
			args: args{
				email:    "some@email.com",
				password: "123123",
			},
			userInDB: nil,
			errFromRepositoryOnFind: &userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			},
			want: false,
			wantErr: &StorageError{storageErrorText: (&userRepository.CommonRepositoryError{
				Action:    "find",
				ErrorText: "test",
			}).Error()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup (before each)
			urCtrl, userRepositoryMock, _, tokenProviderMock := prepareMocks(t)

			u := &userServiceImpl{
				tokenProvider:  tokenProviderMock,
				userRepository: userRepositoryMock,
				tokenTTL:       time.Duration(10000),
			}

			userRepositoryMock.
				EXPECT().
				FindOne(userRepository.FindCondition{Email: tt.args.email}).
				Return(tt.userInDB, tt.errFromRepositoryOnFind)

			got, err := u.ValidateCredentials(tt.args.email, tt.args.password)

			assert.Equal(t, tt.want, got)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			urCtrl.Finish()
		})
	}
}

func prepareMocks(t *testing.T) (
	userRepositoryCtrl *gomock.Controller,
	userRepositoryMock *mocks.MockUserRepository,
	tokenProviderCtrl *gomock.Controller,
	tokenProviderMock *mocks.MockTokenProvider,
) {
	userRepositoryCtrl = gomock.NewController(t)
	userRepositoryMock = mocks.NewMockUserRepository(userRepositoryCtrl)

	tokenProviderCtrl = gomock.NewController(t)
	tokenProviderMock = mocks.NewMockTokenProvider(tokenProviderCtrl)

	return userRepositoryCtrl, userRepositoryMock, tokenProviderCtrl, tokenProviderMock
}

func assertUserDTOAndError(t *testing.T, want entity.UserDTO, wantErr error, got entity.UserDTO, gotErr error) {
	if want == nil {
		assert.Nil(t, got)
	} else {
		assert.NotNil(t, got)
		if got != nil {
			assert.Equal(t, want.ID(), got.ID())
			assert.Equal(t, want.Email(), got.Email())
		}
	}

	if wantErr == nil {
		assert.Nil(t, gotErr)
	} else {
		assert.EqualError(t, gotErr, wantErr.Error())
	}
}
