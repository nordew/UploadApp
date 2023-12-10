package psqldb

import (
	"context"
	"database/sql"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

// AuthStorage defines the interface for storing and retrieving refresh tokens in the database.
type AuthStorage interface {
	// Get retrieves a refresh token session from the database based on the provided token string.
	// Returns the refresh token session if found, or an error if not found or an error occurred.
	Get(ctx context.Context, token string) (*entity.RefreshTokenSession, error)

	// Create stores a new refresh token session in the database.
	// Returns an error if the storage operation fails.
	Create(ctx context.Context, token *entity.RefreshTokenSession) error
}

type authStorage struct {
	db *sql.DB
}

func NewAuthStorage(db *sql.DB) *authStorage {
	return &authStorage{db}
}

func (s *authStorage) Create(ctx context.Context, token *entity.RefreshTokenSession) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user_id=$1", token.UserID)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token, role, expires_at) VALUES($1, $2, $3, $4)`,
		&token.UserID, &token.RefreshToken, &token.Role, &token.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *authStorage) Get(ctx context.Context, token string) (*entity.RefreshTokenSession, error) {
	var refreshToken entity.RefreshTokenSession

	row := s.db.QueryRowContext(ctx,
		"SELECT user_id, token, role,  expires_at FROM refresh_tokens WHERE refresh_token=$1",
		token)

	if err := row.Scan(&refreshToken.UserID, &refreshToken.RefreshToken, &refreshToken.Role, &refreshToken.ExpiresAt); err != nil {
		return nil, err
	}

	return &refreshToken, nil
}
