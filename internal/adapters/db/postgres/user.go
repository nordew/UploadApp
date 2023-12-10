package psqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nordew/UploadApp/internal/domain/entity"
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

	// GetByCredentials retrieves a user from the database by email.
	// It returns an error if the operation fails or the user is not found.
	GetByCredentials(ctx context.Context, email string) (*entity.User, error)

	// CreateRefreshToken updates the refresh token for a user with the specified ID.
	CreateRefreshToken(ctx context.Context, token string, id string) error
}

type userStorage struct {
	db *sql.DB
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
	userId, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, password, registered_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userId, user.Name, user.Email, user.Password, user.RegisteredAt)

	if err != nil {
		if IsDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateKey, err)
		}

		return fmt.Errorf("%w: %v", ErrFailedToInsert, err)
	}

	return nil
}

func (s *userStorage) GetByCredentials(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, email, password, photos_uploaded, role
		FROM users
		WHERE email = $1`,
		email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.PhotosUploaded, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found for email %s", ErrUserNotFound, email)
		}
		return nil, fmt.Errorf("%w: %v", ErrFailedToDecode, err)
	}

	return &user, nil
}

func (s *userStorage) CreateRefreshToken(ctx context.Context, token string, id string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE users SET refresh_token = $1 WHERE id = $2;", token, id)
	if err != nil {
		return err
	}

	return nil
}
