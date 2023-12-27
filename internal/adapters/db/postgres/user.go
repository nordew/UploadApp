package psqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/sirupsen/logrus"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrFailedToMarshal = errors.New("failed to marshal user")
	ErrFailedToInsert  = errors.New("failed to insert user")
	ErrFailedToDecode  = errors.New("failed to decode user")
	ErrDuplicateKey    = errors.New("user with this email already exists")
)

// UserStorage is an interface for user data storage operations.
type UserStorage interface {
	// Create creates a new user in the database.
	Create(ctx context.Context, user entity.User) error

	// GetByCredentials retrieves a user from the database based on either email or ID.
	// It returns an error if the operation fails or the user is not found.
	// The 'identifier' parameter is used for both email and ID, and 'byEmail' indicates the type of identifier.
	GetByCredentials(ctx context.Context, identifier string, byEmail bool) (*entity.User, error)

	// CreateRefreshToken updates the refresh token for a user with the specified ID.
	CreateRefreshToken(ctx context.Context, token, id string) error

	// RefreshSession updates the refresh token for a user with the specified oldToken to a newToken.
	// It returns an error if the operation fails or the oldToken is not found.
	RefreshSession(ctx context.Context, oldToken, newToken string) error

	// ChangePassword changes the user password
	// It returns an error if the operation fails or password stored in db doesn't match with old one.
	ChangePassword(ctx context.Context, email, old, new string) error

	// IncrementPhotosUploaded increments photos_uploaded field in database
	// It returns an error if the operation fails
	IncrementPhotosUploaded(ctx context.Context, userId string) error
}

type userStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewUserStorage(db *sql.DB) *userStorage {
	return &userStorage{
		db: db,
	}
}

func IsDuplicateKeyError(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code == "23505"
}

func (s *userStorage) Create(ctx context.Context, user entity.User) error {
	logger := s.logger.WithField("function", "Create")

	userId, err := uuid.NewUUID()
	if err != nil {
		logger.WithError(err).Error("failed to generate UUID")
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, password, registered_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userId, user.Name, user.Email, user.Password, user.RegisteredAt)

	if err != nil {
		if IsDuplicateKeyError(err) {
			logger.WithError(err).Error("duplicate key error")
			return fmt.Errorf("%w: %v", ErrDuplicateKey, err)
		}

		logger.WithError(err).Error("failed to insert user")
		return fmt.Errorf("%w: %v", ErrFailedToInsert, err)
	}

	return nil
}

func (s *userStorage) GetByCredentials(ctx context.Context, identifier string, byEmail bool) (*entity.User, error) {
	logger := s.logger.WithField("function", "GetByCredentials")

	var user entity.User

	var query string
	var args []interface{}

	if byEmail {
		query = `
			SELECT id, name, email, password, photos_uploaded, role
			FROM users
			WHERE email = $1
		`
		args = append(args, identifier)
	} else {
		query = `
			SELECT id, name, email, password, photos_uploaded, role
			FROM users
			WHERE id = $1
		`

		args = append(args, identifier)
	}

	row := s.db.QueryRowContext(ctx, query, args...)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.PhotosUploaded, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WithError(err).Errorf("user not found for identifier %s", identifier)
			return nil, fmt.Errorf("%w: user not found for identifier %s", ErrUserNotFound, identifier)
		}
		logger.WithError(err).Error("failed to decode user")
		return nil, fmt.Errorf("%w: %v", ErrFailedToDecode, err)
	}

	return &user, nil
}

func (s *userStorage) CreateRefreshToken(ctx context.Context, token string, id string) error {
	logger := s.logger.WithField("function", "CreateRefreshToken")

	_, err := s.db.ExecContext(ctx, "UPDATE users SET refresh_token = $1 WHERE id = $2;", token, id)
	if err != nil {
		logger.WithError(err).Error("failed to create refresh token")
		return err
	}

	return nil
}

func (s *userStorage) RefreshSession(ctx context.Context, oldToken string, newToken string) error {
	logger := s.logger.WithField("function", "RefreshSession")

	_, err := s.db.ExecContext(ctx, "UPDATE users SET refresh_token = $1 WHERE refresh_token = $2;", oldToken, newToken)
	if err != nil {
		logger.WithError(err).Error("no such refresh token")
		return fmt.Errorf("no such refresh token")
	}

	return nil
}

func (s *userStorage) ChangePassword(ctx context.Context, email, old, new string) error {
	logger := s.logger.WithField("function", "ChangePassword")

	var dbPassword string

	row := s.db.QueryRowContext(ctx, "SELECT password FROM users WHERE email = $1", email)

	if err := row.Scan(&dbPassword); err != nil {
		logger.WithError(err).Error("failed to get password from database")
		return err
	}

	if dbPassword != old {
		logger.Error("invalid old password")
		return fmt.Errorf("invalid old password")
	}

	_, err := s.db.ExecContext(ctx, "UPDATE users SET password = $1 WHERE email = $2", new, email)
	if err != nil {
		logger.WithError(err).Error("failed to change password")
		return err
	}

	return nil
}

func (s *userStorage) IncrementPhotosUploaded(ctx context.Context, userId string) error {
	logger := s.logger.WithField("function", "IncrementPhotosUploaded")

	_, err := s.db.ExecContext(ctx, "UPDATE users SET photos_uploaded = photos_uploaded + 1 WHERE id = $1;", userId)
	if err != nil {
		logger.WithError(err).Error("failed to increment photos uploaded count")
		return err
	}

	return nil
}
