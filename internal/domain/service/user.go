package service

import (
	"context"

	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
)

type Users interface {
	// SignUp creates a new user account.
	SignUp(ctx context.Context, input entity.SignUpInput) error
	// SignIn retrieves a user from the database by email and password.
	// It returns an access token and an error if the operation fails or the user is not found.
	SignIn(ctx context.Context, input entity.SignInInput) (string, error)
}

type UserService struct {
	storage mongodb.UserStorage
	hasher  hasher.PasswordHasher
	token   auth.JWTToken
}

func NewUserService(storage mongodb.UserStorage, hasher hasher.PasswordHasher, token auth.JWTToken) *UserService {
	return &UserService{
		storage: storage,
		hasher:  hasher,
		token:   token,
	}
}

func (s *UserService) SignUp(ctx context.Context, input entity.SignUpInput) error {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	if err := input.Validate(); err != nil {
		return err
	}

	user := entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	return s.storage.Create(context.TODO(), &user)
}

func (s *UserService) SignIn(ctx context.Context, input entity.SignInInput) (string, error) {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	if err := input.Validate(); err != nil {
		return "", err
	}

	user, err := s.storage.GetByCredentials(context.TODO(), input.Email, hashedPassword)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	accessToken, err := s.token.NewJWTToken(user.ID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
