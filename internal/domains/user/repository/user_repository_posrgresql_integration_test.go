// +build integration

package repository

import (
	"github.com/Viverov/guideliner/internal/config"
	"github.com/Viverov/guideliner/internal/db"
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
	type userData struct {
		id uint
		email string
		password string
	}
	tests := []struct {
		before func()
		name    string
		args    args
		wantUser bool
		wantUserData userData
		wantErr error
		after func()
	}{
		{
			before: func () {
				dbInstance.Create(&userModel{
					Model: gorm.Model{ID: 10},
					Email:    "email@test.com",
					Password: "pass",
				})
			},
			name:    "Should return existing record by ID",
			args:    args{condition: FindCondition{
				ID:    10,
			}},
			wantUser: true,
			wantUserData: userData{
				id:       10,
				email: 	  "email@test.com",
				password: "pass",
			},
			wantErr: nil,
			after: func() {
				dbInstance.Delete(&userModel{
					Model:    gorm.Model{ID: 10},
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
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
				}
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			tt.after()
		})
	}
}
//
//func Test_userRepositoryPostgresql_Insert(t *testing.T) {
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
//		wantId  uint
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &userRepositoryPostgresql{
//				db: tt.fields.db,
//			}
//			gotId, err := r.Insert(tt.args.u)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if gotId != tt.wantId {
//				t.Errorf("Insert() gotId = %v, want %v", gotId, tt.wantId)
//			}
//		})
//	}
//}
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
