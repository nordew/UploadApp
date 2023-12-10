package service

import (
	"context"
	"fmt"
	psqldb "github.com/nordew/UploadApp/internal/adapters/db/postgres"
	"time"

	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
	"github.com/pkg/errors"
)

var (
	ErrValidationFailed = errors.New("invlid input")
)

// Users is the interface that defines methods for user-related operations, such as sign-up and sign-in.
type Users interface {
	// SignUp creates a new user account based on the provided sign-up input.
	// It returns an error if the operation fails, including cases where the provided email already exists.
	SignUp(ctx context.Context, input entity.SignUpInput) error

	// SignIn retrieves a user from the database by email and password, and generates an access token.
	// It returns the generated access token and an error if the operation fails or the user is not found.
	SignIn(ctx context.Context, input entity.SignInInput) (string, string, error)

	Refresh(ctx context.Context, refreshToken string) (string, error)
}

type UserService struct {
	authService Auths
	storage     psqldb.UserStorage
	hasher      hasher.PasswordHasher
	auth        auth.Authenticator

	hmacSecret string
}

func NewUserService(authService Auths, storage psqldb.UserStorage, hasher hasher.PasswordHasher, auth auth.Authenticator, hmacSecret string) *UserService {
	return &UserService{
		authService: authService,
		storage:     storage,
		hasher:      hasher,
		auth:        auth,
		hmacSecret:  hmacSecret,
	}
}

func (s *UserService) SignUp(ctx context.Context, input entity.SignUpInput) error {
	if err := input.Validate(); err != nil {
		return ErrValidationFailed
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	user := entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := s.storage.Create(ctx, user); err != nil {
		switch {
		case errors.Is(err, psqldb.ErrDuplicateKey):
			return fmt.Errorf("email already exists: %w", err)
		case errors.Is(err, psqldb.ErrFailedToInsert):
			return fmt.Errorf("failed to create user: %w", err)
		default:
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	return nil
}

func (s *UserService) SignIn(ctx context.Context, input entity.SignInInput) (string, string, error) {
	if err := input.Validate(); err != nil {
		return "", "", fmt.Errorf("invalid input: %w", err)
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.storage.GetByCredentials(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, psqldb.ErrUserNotFound):
			return "", "", fmt.Errorf("authentication failed: %w", err)
		case errors.Is(err, psqldb.ErrFailedToDecode):
			return "", "", fmt.Errorf("authentication failed: %w", err)
		default:
			return "", "", fmt.Errorf("authentication failed: unexpected error: %w", err)
		}
	}

	if hashedPassword != user.Password {
		return "", "", fmt.Errorf("wrong password")
	}

	accessToken, refreshToken, err := s.auth.GenerateTokens(&auth.GenerateTokenClaimsOptions{
		UserId: user.ID,
		Role:   user.Role,
	})
	if err != nil {
		return "", "", err
	}

	if err := s.authService.Create(ctx, &entity.RefreshTokenSession{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		Role:         user.Role,
		ExpiresAt:    time.Now().Add(24 * 30 * time.Hour),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}
