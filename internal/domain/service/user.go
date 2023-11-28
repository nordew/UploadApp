package service

import (
	"context"
	"fmt"

	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
	"github.com/pkg/errors"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
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
	auth    auth.Authenticator

	hmacSecret string
}

func NewUserService(storage mongodb.UserStorage, hasher hasher.PasswordHasher, auth auth.Authenticator, hmacSecret string) *UserService {
	return &UserService{
		storage:    storage,
		hasher:     hasher,
		auth:       auth,
		hmacSecret: hmacSecret,
	}
}

func (s *UserService) SignUp(ctx context.Context, input entity.SignUpInput) error {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "invalid input")
	}

	user := entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := s.storage.Create(ctx, user); err != nil {
		switch {
		case errors.Is(err, mongodb.ErrDuplicateKey):
			return fmt.Errorf("email already exists: %w", err)
		case errors.Is(err, mongodb.ErrFailedToInsert):
			return fmt.Errorf("failed to create user: %w", err)
		default:
			return fmt.Errorf("unexpected error: %w", err)
		}
	}

	return nil
}

func (s *UserService) SignIn(ctx context.Context, input entity.SignInInput) (string, error) {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	if err := input.Validate(); err != nil {
		return "", err
	}

	user, err := s.storage.GetByCredentials(ctx, input.Email, hashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, mongodb.ErrUserNotFound):
			return "", fmt.Errorf("authentication failed: %w", err)
		case errors.Is(err, mongodb.ErrFailedToDecode):
			return "", fmt.Errorf("authentication failed: %w", err)
		default:
			return "", fmt.Errorf("authentication failed: unexpected error: %w", err)
		}
	}

	if err != nil {
		return "", err
	}

	accessToken, err := s.auth.GenerateToken(&auth.GenerateTokenClaimsOptions{
		UserId:   user.ID,
		UserName: user.Name,
	})
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
