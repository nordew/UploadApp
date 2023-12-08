package service

import (
	"context"
	"fmt"
	psqldb "github.com/nordew/UploadApp/internal/adapters/db/postgres"

	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
	"github.com/pkg/errors"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrValidationFailed   = errors.New("invlid input")
)

// Users is the interface that defines methods for user-related operations, such as sign-up and sign-in.
type Users interface {
	// SignUp creates a new user account based on the provided sign-up input.
	// It returns an error if the operation fails, including cases where the provided email already exists.
	SignUp(ctx context.Context, input entity.SignUpInput) (string, error)

	// SignIn retrieves a user from the database by email and password, and generates an access token.
	// It returns the generated access token and an error if the operation fails or the user is not found.
	SignIn(ctx context.Context, input entity.SignInInput) (string, error)
}

type UserService struct {
	storage psqldb.UserStorage
	hasher  hasher.PasswordHasher
	auth    auth.Authenticator

	hmacSecret string
}

func NewUserService(storage psqldb.UserStorage, hasher hasher.PasswordHasher, auth auth.Authenticator, hmacSecret string) *UserService {
	return &UserService{
		storage:    storage,
		hasher:     hasher,
		auth:       auth,
		hmacSecret: hmacSecret,
	}
}

func (s *UserService) SignUp(ctx context.Context, input entity.SignUpInput) (string, error) {
	if err := input.Validate(); err != nil {
		return "", ErrValidationFailed
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	user := entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	id, err := s.storage.Create(ctx, user)

	if err != nil {
		switch {
		case errors.Is(err, psqldb.ErrDuplicateKey):
			return "", fmt.Errorf("email already exists: %w", err)
		case errors.Is(err, psqldb.ErrFailedToInsert):
			return "", fmt.Errorf("failed to create user: %w", err)
		default:
			return "", fmt.Errorf("failed to create user: %w", err)
		}
	}

	return id, nil
}

func (s *UserService) SignIn(ctx context.Context, input entity.SignInInput) (string, error) {
	if err := input.Validate(); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.storage.GetByCredentials(ctx, input.Email, hashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, psqldb.ErrUserNotFound):
			return "", fmt.Errorf("authentication failed: %w", err)
		case errors.Is(err, psqldb.ErrFailedToDecode):
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
