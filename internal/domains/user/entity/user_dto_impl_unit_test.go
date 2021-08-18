// +build unit

package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUserDTOFromEntity(t *testing.T) {
	type args struct {
		user User
	}
	tests := []struct {
		name string
		args args
		want *userDTOImpl
	}{
		{
			name: "Should create DTO",
			args: args{user: &userImpl{
				id:       10,
				email:    "example@email.com",
				password: "supersecret",
			}},
			want: &userDTOImpl{
				id:    10,
				email: "example@email.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserDTOFromEntity(tt.args.user)
			assert.Equal(t, tt.want.ID(), got.ID())
			assert.Equal(t, tt.want.Email(), got.Email())
		})
	}
}

func Test_userDTOImpl_Email(t *testing.T) {
	type fields struct {
		id    uint
		email string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return email",
			fields: fields{
				id:    10,
				email: "example@email.com",
			},
			want: "example@email.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			D := &userDTOImpl{
				id:    tt.fields.id,
				email: tt.fields.email,
			}
			assert.Equal(t, tt.want, D.Email())
		})
	}
}

func Test_userDTOImpl_ID(t *testing.T) {
	type fields struct {
		id    uint
		email string
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "Should return ID",
			fields: fields{
				id:    10,
				email: "test@test.com",
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			D := &userDTOImpl{
				id:    tt.fields.id,
				email: tt.fields.email,
			}
			assert.Equal(t, tt.want, D.ID())
		})
	}
}
