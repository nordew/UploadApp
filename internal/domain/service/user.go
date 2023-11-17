package service

import (
	"context"

	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

type Users interface {
	SignUp(ctx context.Context, input entity.SignUpInput) error
	//SignIn()
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserService struct {
	repo   mongodb.UserStorage
	hasher PasswordHasher
}

func NewUserService(repo mongodb.UserStorage, hasher PasswordHasher) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
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
