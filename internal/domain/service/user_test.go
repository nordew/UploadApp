package service

import (
	"context"
	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/internal/mocks"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		storage    mongodb.UserStorage
		hasher     hasher.PasswordHasher
		auth       auth.Authenticator
		hmacSecret string
	}
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.storage, tt.args.hasher, tt.args.auth, tt.args.hmacSecret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_SignIn(t *testing.T) {

	type args struct {
		ctx   context.Context
		input entity.SignInInput
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Roman SignIn",
			args: args{
				ctx: context.Background(),
				input: entity.SignInInput{
					Email:    "roman@gmail.com",
					Password: "roman123",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		mockStorage := mocks.NewUserStorage(t)
		mockHasher := mocks.NewPasswordHasher(t)
		mockAuth := mocks.NewAuthenticator(t)

		mockStorage.
			On("GetByCredentials", tt.args.ctx, tt.args.input.Email, tt.args.input.Password).
			Return(&entity.User{}, nil)

		mockHasher.
			On("Hash", tt.args.input.Password).
			Return("roman123", nil)

		mockAuth.On("GenerateToken", auth.GenerateTokenClaimsOptions{})

		t.Run(tt.name, func(t *testing.T) {
			s := &UserService{
				storage:    mockStorage,
				hasher:     mockHasher,
				auth:       mockAuth,
				hmacSecret: "norman!28",
			}

			_, err := s.SignIn(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserService_SignUp(t *testing.T) {
	type args struct {
		ctx   context.Context
		input entity.SignUpInput
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Test_1",
			args: args{
				ctx: context.Background(),
				input: entity.SignUpInput{
					"John",
					"John.doe@gmail.com",
					"john123",
				},
			},
			wantErr: nil,
		},
		{
			name: "Test_2",
			args: args{
				ctx: context.Background(),
				input: entity.SignUpInput{
					"Jane",
					"Jane.smith@gmail.com",
					"jane456",
				},
			},
			wantErr: nil,
		},
		{
			name: "Test_3_InvalidInput",
			args: args{
				ctx: context.Background(),
				input: entity.SignUpInput{
					"Yan",
					"invalid_email",
					"password123",
				},
			},
			wantErr: ErrValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewUserStorage(t)
			mockHasher := mocks.NewPasswordHasher(t)
			mockAuth := mocks.NewAuthenticator(t)

			mockHasher.
				On("Hash", tt.args.input.Password).
				Return(tt.args.input.Password, nil)

			mockStorage.
				On("Create", tt.args.ctx, entity.User{
					Name:     tt.args.input.Name,
					Email:    tt.args.input.Email,
					Password: tt.args.input.Password,
				}).
				Return(nil)

			s := &UserService{
				storage:    mockStorage,
				hasher:     mockHasher,
				auth:       mockAuth,
				hmacSecret: "",
			}

			err := s.SignUp(tt.args.ctx, tt.args.input)
			if err != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
			assert.Equal(t, err, tt.wantErr)
		})
	}

}
