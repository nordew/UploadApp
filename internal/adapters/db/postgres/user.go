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

type UserStorage interface {
	// Create creates a new user in the database.
	Create(ctx context.Context, user entity.User) (string, error)
	// GetByCredentials retrieves a user from the database by email and password.
	// It returns an error if the operation fails or the user is not found.
	GetByCredentials(ctx context.Context, email, password string) (*entity.User, error)
}

type userStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *userStorage {
	return &userStorage{db}
}

func IsDuplicateKeyError(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code == "23505"
}

func (s *userStorage) Create(ctx context.Context, user entity.User) (string, error) {
	userId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, password, registered_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userId, user.Name, user.Email, user.Password, user.RegisteredAt)

	if err != nil {
		if IsDuplicateKeyError(err) {
			return "", fmt.Errorf("%w: %v", ErrDuplicateKey, err)
		}

		return "", fmt.Errorf("%w: %v", ErrFailedToInsert, err)
	}

	return userId.String(), nil
}

func (s *userStorage) GetByCredentials(ctx context.Context, email, password string) (*entity.User, error) {
	var user entity.User

	row := s.db.QueryRowContext(ctx, `
		SELECT id, email, password
		FROM users
		WHERE email = $1 AND password = $2`,
		email, password)

	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user not found for email %s", ErrUserNotFound, email)
		}
		return nil, fmt.Errorf("%w: %v", ErrFailedToDecode, err)
	}

	return &user, nil
}
