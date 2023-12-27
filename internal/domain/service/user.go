package service

import (
	"context"
	"fmt"
	psqldb "github.com/nordew/UploadApp/internal/adapters/db/postgres"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	ErrValidationFailed = errors.New("invalid input")
)

// Users is the interface that defines methods for user-related operations, such as sign-up and sign-in.
type Users interface {
	// SignUp creates a new user account based on the provided sign-up input.
	// It returns an error if the operation fails, including cases where the provided email already exists.
	SignUp(ctx context.Context, input entity.SignUpInput) error

	// SignIn retrieves a user from the database by email and password, and generates an access token.
	// It returns the generated access token and an error if the operation fails or the user is not found.
	SignIn(ctx context.Context, input entity.SignInInput) (string, string, error)

	// Refresh generates new access and refresh tokens for an existing user session based on the provided ID and role.
	// It returns the new access and refresh tokens, and an error if the operation fails.
	Refresh(ctx context.Context, id, role string) (string, string, error)

	GetCredentials(ctx context.Context, identifier string, byEmail bool) (*entity.User, error)

	// ChangePassword updates the user's password in the system.
	// hashes the old and new passwords, and updates the password in the storage.
	// If the operation is successful, it returns nil. Otherwise, it returns an error.
	// The error may indicate validation failures, hashing errors, or storage-related issues
	ChangePassword(ctx context.Context, id, old, new string) error

	IncrementPhotosUploaded(ctx context.Context, id string) error
}

type UserService struct {
	storage psqldb.UserStorage
	hasher  hasher.PasswordHasher
	auth    auth.Authenticator
	logger  *logrus.Logger

	hmacSecret string
}

func NewUserService(storage psqldb.UserStorage, hasher hasher.PasswordHasher, auth auth.Authenticator, logger *logrus.Logger, hmacSecret string) *UserService {
	return &UserService{
		storage:    storage,
		hasher:     hasher,
		auth:       auth,
		logger:     logger,
		hmacSecret: hmacSecret,
	}
}

func (s *UserService) SignUp(ctx context.Context, input entity.SignUpInput) error {
	if err := input.Validate(); err != nil {
		s.logger.WithError(err).Error("SignUp: validation failed")
		return ErrValidationFailed
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		s.logger.WithError(err).Error("SignUp: failed to hash password")
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
			s.logger.WithError(err).Error("SignUp: email already exists")
			return fmt.Errorf("SignUp: email already exists: %w", err)
		case errors.Is(err, psqldb.ErrFailedToInsert):
			s.logger.WithError(err).Error("SignUp: failed to create user")
			return fmt.Errorf("SignUp: failed to create user: %w", err)
		default:
			s.logger.WithError(err).Error("SignUp: failed to create user")
			return fmt.Errorf("SignUp: failed to create user: %w", err)
		}
	}

	s.logger.Info("SignUp: user created successfully")
	return nil
}

func (s *UserService) SignIn(ctx context.Context, input entity.SignInInput) (string, string, error) {
	if err := input.Validate(); err != nil {
		return "", "", fmt.Errorf("invalid input: %w", err)
	}

	hashedPasswordCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		hashedPassword, err := s.hasher.Hash(input.Password)
		if err != nil {
			s.logger.WithError(err).Error("failed to hash password")
			errCh <- err
			return
		}
		hashedPasswordCh <- hashedPassword
	}()

	userCh := make(chan *entity.User, 1)
	go func() {
		user, err := s.storage.GetByCredentials(ctx, input.Email, true)
		if err != nil {
			errCh <- err
			return
		}
		userCh <- user
	}()

	var hashedPassword string
	var user *entity.User
	for i := 0; i < 2; i++ {
		select {
		case err := <-errCh:
			return "", "", fmt.Errorf("authentication failed: %w", err)
		case hashedPassword = <-hashedPasswordCh:
		case user = <-userCh:
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

	if err := s.storage.CreateRefreshToken(context.Background(), refreshToken, user.ID); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) Refresh(ctx context.Context, id, role string) (string, string, error) {
	var accessToken, refreshToken string
	var genErr, refreshErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		var err error
		accessToken, refreshToken, err = s.auth.GenerateTokens(&auth.GenerateTokenClaimsOptions{
			UserId: id,
			Role:   role,
		})
		if err != nil {
			genErr = err
		}
	}()

	go func() {
		defer wg.Done()
		refreshErr = s.storage.RefreshSession(ctx, refreshToken, id)
	}()

	wg.Wait()

	if genErr != nil {
		return "", "", genErr
	}

	if refreshErr != nil {
		return "", "", refreshErr
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) GetCredentials(ctx context.Context, identifier string, byEmail bool) (*entity.User, error) {
	return s.storage.GetByCredentials(ctx, identifier, byEmail)
}

func (s *UserService) ChangePassword(ctx context.Context, email, old, new string) error {
	hashedOldPassword, err := s.hasher.Hash(old)
	if err != nil {
		return fmt.Errorf("failed to hash old paasword")
	}

	hashedNewPassword, err := s.hasher.Hash(new)
	if err != nil {
		return fmt.Errorf("failed to hash new paasword")
	}

	return s.storage.ChangePassword(ctx, email, hashedOldPassword, hashedNewPassword)
}

func (s *UserService) IncrementPhotosUploaded(ctx context.Context, id string) error {
	return s.storage.IncrementPhotosUploaded(ctx, id)
}
