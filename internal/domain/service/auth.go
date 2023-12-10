package service

import (
	"context"
	"fmt"
	psqldb "github.com/nordew/UploadApp/internal/adapters/db/postgres"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"time"
)

// Auths defines the authentication service interface, which includes methods for handling refresh tokens and creating new tokens.
type Auths interface {
	// RefreshTokens generates new access and refresh tokens based on the provided refresh token and user information.
	// It checks the validity of the refresh token and returns the new access and refresh tokens.
	// Returns an error if the refresh token is invalid or expired.
	RefreshTokens(ctx context.Context, token string) (string, string, error)

	// Create stores the given refresh token session in the storage.
	// This method is used during the initial creation of refresh tokens.
	Create(ctx context.Context, token *entity.RefreshTokenSession) error
}

type authService struct {
	authStorage psqldb.AuthStorage
	auth        auth.Authenticator
}

func NewAuthService(authStorage psqldb.AuthStorage) *authService {
	return &authService{authStorage: authStorage}
}

func (s *authService) Create(ctx context.Context, token *entity.RefreshTokenSession) error {
	return s.authStorage.Create(ctx, token)
}

func (s *authService) RefreshTokens(ctx context.Context, token string) (string, string, error) {
	responseRefreshToken, err := s.authStorage.Get(ctx, token)
	if err != nil {
		return "", "", err
	}

	if responseRefreshToken.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", fmt.Errorf("refresh token is expired")
	}

	accessToken, refreshToken, err := s.auth.GenerateTokens(&auth.GenerateTokenClaimsOptions{
		UserId: responseRefreshToken.UserID,
		Role:   responseRefreshToken.Role,
	})

	return accessToken, refreshToken, nil
}
