package service

import (
	"context"

	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
)

type Users interface {
	SignUp(ctx context.Context, input entity.SignUpInput) error
	SignIn(ctx context.Context, input entity.SignInInput) (string, error)
}

type UserService struct {
	repo   mongodb.UserStorage
	hasher hasher.PasswordHasher
	token  auth.JWTToken
}

func NewUserService(repo mongodb.UserStorage, hasher hasher.PasswordHasher, token auth.JWTToken) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
		token:  token,
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

	return s.repo.Create(context.TODO(), &user)
}

func (s *UserService) SignIn(ctx context.Context, input entity.SignInInput) (string, error) {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	if err := input.Validate(); err != nil {
		return "", err
	}

	user, err := s.repo.GetByCredentials(context.TODO(), input.Email, hashedPassword)
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
