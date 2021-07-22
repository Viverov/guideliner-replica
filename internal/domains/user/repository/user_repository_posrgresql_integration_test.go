// +build integration

package repository

import (
	"github.com/Viverov/guideliner/internal/config"
	"github.com/Viverov/guideliner/internal/db"
	"github.com/Viverov/guideliner/internal/domains/user/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var cfg = config.InitConfig(config.EnvTest, "./config.json")
var dbInstance = db.GetDB(&db.DBOptions{
	Host:     cfg.DB.Host,
	Port:     cfg.DB.Port,
	Login:    cfg.DB.Login,
	Password: cfg.DB.Password,
	Name:     cfg.DB.Name,
	SSLMode:  cfg.DB.SSLMode,
})

type userData struct {
	id       uint
	email    string
	password string
}

func TestNewUserRepositoryPostgresql(t *testing.T) {
	dbInstance := &gorm.DB{}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should return new repository",
			args: args{db: dbInstance},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserRepositoryPostgresql(tt.args.db)
			assert.NotNil(t, got)
		})
	}
}

func Test_userRepositoryPostgresql_FindOne(t *testing.T) {
	type args struct {
		condition FindCondition
	}

	// Setup test data (before all)
	testUserData := userData{
		id:       10,
		email:    "test@test.com",
		password: "pass",
	}
	dbInstance.Create(&userModel{
		Model:    gorm.Model{ID: testUserData.id},
		Email:    testUserData.email,
		Password: testUserData.password,
	})

	tests := []struct {
		name         string
		args         args
		wantUser     bool
		wantUserData userData
		wantErr      error
	}{
		{
			name: "Should return existing record by ID",
			args: args{condition: FindCondition{
				ID: testUserData.id,
			}},
			wantUser:     true,
			wantUserData: testUserData,
			wantErr:      nil,
		},
		{
			name: "Should return user record by email",
			args: args{condition: FindCondition{
				Email: testUserData.email,
			}},
			wantUser:     true,
			wantUserData: testUserData,
			wantErr:      nil,
		},
		{
			name: "Should return nil for undefined ID",
			args: args{condition: FindCondition{
				ID: 1234,
			}},
			wantUser:     false,
			wantUserData: userData{},
			wantErr:      nil,
		},
		{
			name: "Should return nil for undefined email",
			args: args{condition: FindCondition{
				Email: "random@email.com",
			}},
			wantUser:     false,
			wantUserData: userData{},
			wantErr:      nil,
		},
		{
			name: "Should return user by email and id",
			args: args{
				condition: FindCondition{
					ID:    testUserData.id,
					Email: testUserData.email,
				},
			},
			wantUser:     true,
			wantUserData: testUserData,
			wantErr:      nil,
		},
		{
			name:         "Should return error for empty args",
			args:         args{condition: FindCondition{}},
			wantUser:     false,
			wantUserData: userData{},
			wantErr:      &InvalidFindConditionError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepositoryPostgresql{
				db: dbInstance,
			}
			got, err := r.FindOne(tt.args.condition)

			if tt.wantErr == nil {
				assert.Nil(t, err)
				if tt.wantUser {
					assert.Equal(t, tt.wantUserData.id, got.ID())
					assert.Equal(t, tt.wantUserData.email, got.Email())
					assert.Equal(t, tt.wantUserData.password, got.Password())
				} else {
					assert.Nil(t, got)
				}
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}

	// Clean up (after all)
	dbInstance.Where("1 = 1").Delete(&userModel{})
}

func Test_userRepositoryPostgresql_Insert(t *testing.T) {
	type args struct {
		u entity.User
	}

	// Setup test data (before all)
	alreadyExistsUserData := userData{
		email:    "already@exists.com",
		password: "already_pass",
	}
	alreadyExistsUser, _ := entity.NewUser(0, alreadyExistsUserData.email, alreadyExistsUserData.password)
	dbInstance.Create(&userModel{
		Email:    alreadyExistsUserData.email,
		Password: alreadyExistsUserData.password,
	})

	testUserData := userData{
		email:    "test1@test.com",
		password: "pass",
	}
	insertedUser, _ := entity.NewUser(0, testUserData.email, testUserData.password)

	tests := []struct {
		name    string
		args    args
		wantId  bool
		wantErr error
	}{
		{
			name:    "Should insert user into database and return ID",
			args:    args{u: insertedUser},
			wantId:  true,
			wantErr: nil,
		},
		{
			name:    "Should return error for already exists user",
			args:    args{u: alreadyExistsUser},
			wantId:  false,
			wantErr: &UserAlreadyExistsError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepositoryPostgresql{
				db: dbInstance,
			}

			gotId, err := r.Insert(tt.args.u)
			if tt.wantId == true {
				assert.Nil(t, err)
				assert.NotEqual(t, 0, tt.wantId)

				// Check creation
				userFromDB := &userModel{}
				result := r.db.Where(&userModel{Model: gorm.Model{ID: gotId}}).First(userFromDB)
				assert.Nil(t, result.Error)
				assert.Equal(t, gotId, userFromDB.ID)
				assert.Equal(t, testUserData.email, userFromDB.Email)
				assert.Equal(t, testUserData.password, userFromDB.Password)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}

			// Clean up (after each)
			if gotId != 0 {
				dbInstance.Where("id = ?", gotId).Delete(&userModel{})
			}
		})
	}

	// Clean up (after all)
	dbInstance.Where("1 = 1").Delete(&userModel{})
}

//
//func Test_userRepositoryPostgresql_Update(t *testing.T) {
//	type fields struct {
//		db *gorm.DB
//	}
//	type args struct {
//		u entity.User
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &userRepositoryPostgresql{
//				db: tt.fields.db,
//			}
//			if err := r.Update(tt.args.u); (err != nil) != tt.wantErr {
//				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
