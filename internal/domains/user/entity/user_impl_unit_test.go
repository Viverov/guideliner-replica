// +build unit

package entity

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr error
	}{
		{
			name:    "Correct",
			args:    args{email: "email", password: "1234567890"},
			want:    &userImpl{email: "email"},
			wantErr: nil,
		},
		{
			name:    "Should lowercase email",
			args:    args{email: "EmAiL@eMaIl.com", password: "1234567890"},
			want:    &userImpl{email: "email@email.com"},
			wantErr: nil,
		},
		{
			name:    "Empty email",
			args:    args{email: "", password: "1234567890"},
			want:    nil,
			wantErr: &EmptyArgError{argName: argNameEmail},
		},
		{
			name:    "Empty Password",
			args:    args{email: "Email", password: ""},
			want:    nil,
			wantErr: &EmptyArgError{argName: argNamePassword},
		},
		{
			name:    "Empty email and password",
			args:    args{email: "", password: ""},
			want:    nil,
			wantErr: &EmptyArgError{argName: argNameEmail},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUser(tt.args.email, tt.args.password)

			if tt.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, strings.ToLower(tt.want.Email()), got.Email())
				assert.True(t, got.ValidatePassword(tt.args.password))
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	type args struct {
		id       uint
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr error
	}{
		{
			name: "Correct",
			args: args{
				id:       0,
				email:    "some@example.com",
				password: "123123123",
			},
			want: &userImpl{
				id:       0,
				email:    "some@example.com",
				password: "123123123",
			},
			wantErr: nil,
		},
		{
			name:    "Empty email",
			args:    args{id: 0, email: "", password: "1234567890"},
			want:    nil,
			wantErr: &EmptyArgError{argName: argNameEmail},
		},
		{
			name:    "Empty Password",
			args:    args{id: 0, email: "Email", password: ""},
			want:    nil,
			wantErr: &EmptyArgError{argName: argNamePassword},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.id, tt.args.email, tt.args.password)

			if tt.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, tt.want.Email(), got.Email())
				assert.Equal(t, tt.want.Password(), got.Password())
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_userImpl_Email(t *testing.T) {
	type fields struct {
		id       uint
		email    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Correct",
			fields: fields{
				id:       0,
				email:    "example@email.com",
				password: "123123",
			},
			want: "example@email.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userImpl{
				id:       tt.fields.id,
				email:    tt.fields.email,
				password: tt.fields.password,
			}
			assert.Equal(t, tt.want, u.Email())
		})
	}
}

func Test_userImpl_ID(t *testing.T) {
	type fields struct {
		id       uint
		email    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			name: "Correct",
			fields: fields{
				id:       500,
				email:    "example@email.com",
				password: "123123",
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userImpl{
				id:       tt.fields.id,
				email:    tt.fields.email,
				password: tt.fields.password,
			}
			assert.Equal(t, tt.want, u.ID())
		})
	}
}

func Test_userImpl_Password(t *testing.T) {
	type fields struct {
		id       uint
		email    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Correct",
			fields: fields{
				id:       0,
				email:    "example@email.com",
				password: "123123",
			},
			want: "123123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userImpl{
				id:       tt.fields.id,
				email:    tt.fields.email,
				password: tt.fields.password,
			}
			assert.Equal(t, tt.want, u.Password())
		})
	}
}

func Test_userImpl_SetID(t *testing.T) {
	type fields struct {
		id       uint
		email    string
		password string
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantID uint
		wantErr error
	}{
		{
			name:    "Should set new id",
			fields:  fields{
				id:       0,
				email:    "",
				password: "",
			},
			args:    args{id: 10},
			wantID: 10,
			wantErr: nil,
		},
		{
			name:    "Should return new error for zero id",
			fields:  fields{
				id:       50,
				email:    "",
				password: "",
			},
			args:    args{id: 0},
			wantID: 50,
			wantErr: &InvalidIdError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userImpl{
				id:       tt.fields.id,
				email:    tt.fields.email,
				password: tt.fields.password,
			}
			err := u.SetID(tt.args.id)
			if tt.wantErr == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			assert.Equal(t, tt.wantID, u.ID())
		})
	}
}

func Test_userImpl_SetEmail(t *testing.T) {
	type fields struct {
		id       uint
		email    string
		password string
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "Correct",
			fields: fields{
				id:       10,
				email:    "example@email.com",
				password: "123123",
			},
			args:    args{email: "new_email@email.com"},
			wantErr: nil,
		},
		{
			name: "Should lowerCase Email",
			fields: fields{
				id:       10,
				email:    "example@email.com",
				password: "123123",
			},
			args:    args{email: "NeW_eMaIl@emAAAAAAAAAAil.coOOOOOOm"},
			wantErr: nil,
		},
		{
			name: "Empty email",
			fields: fields{
				id:       10,
				email:    "example@email.com",
				password: "123123",
			},
			args:    args{email: ""},
			wantErr: &EmptyArgError{argName: argNameEmail},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userImpl{
				id:       tt.fields.id,
				email:    tt.fields.email,
				password: tt.fields.password,
			}
			err := u.SetEmail(tt.args.email)
			if tt.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, strings.ToLower(tt.args.email), u.Email())
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_userImpl_SetPassword(t *testing.T) {
	type fields struct {
		id       uint
		email    string
		password string
	}
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "Correct",
			fields: fields{
				id:       10,
				email:    "email@email.com",
				password: "123456",
			},
			args: args{
				password: "asddfgdfhgfgjh",
			},
			wantErr: nil,
		},
		{
			name: "Empty password",
			fields: fields{
				id:       10,
				email:    "email@email.com",
				password: "123456",
			},
			args: args{
				password: "",
			},
			wantErr: &EmptyArgError{argName: argNamePassword},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userImpl{
				id:       tt.fields.id,
				email:    tt.fields.email,
				password: tt.fields.password,
			}
			err := u.SetPassword(tt.args.password)

			if tt.wantErr == nil {
				assert.Nil(t, err)
				assert.True(t, u.ValidatePassword(tt.args.password))
				assert.False(t, u.ValidatePassword(tt.fields.password))
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func Test_userImpl_ValidatePassword(t *testing.T) {
	type fields struct {
		id       uint
		email    string
		password string
	}
	type args struct {
		password string
	}
	tests := []struct {
		name           string
		password       string
		passedPassword string
		wantIsValid    bool
	}{
		{
			name:           "Correct",
			password:       "12345",
			passedPassword: "12345",
			wantIsValid:    true,
		},
		{
			name:           "Invalid password",
			password:       "12345",
			passedPassword: "qweert",
			wantIsValid:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userImpl{
				id:    10,
				email: "example@email.com",
			}
			// Set & hash password
			_ = u.SetPassword(tt.password)
			assert.Equal(t, tt.wantIsValid, u.ValidatePassword(tt.passedPassword))
		})
	}
}
