package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/hasher"
)

type Users interface {
	// SignUp creates a new user account.
	SignUp(ctx context.Context, input entity.SignUpInput) error

	// SignIn retrieves a user from the database by email and password.
	// It returns an access token and an error if the operation fails or the user is not found.
	SignIn(ctx context.Context, input entity.SignInInput) (string, error)

	// Parsing Token from request
	ParseToken(ctx context.Context, token string) (jwt.MapClaims, error)
}

type UserService struct {
	storage mongodb.UserStorage
	hasher  hasher.PasswordHasher
	auth    auth.Token

	hmacSecret string
}

func NewUserService(storage mongodb.UserStorage, hasher hasher.PasswordHasher, auth auth.Token, hmacSecret string) *UserService {
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

	return s.storage.Create(context.TODO(), user)
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
		return "", err
	}

	if err != nil {
		return "", err
	}

	accessToken, err := s.auth.NewToken(user.ID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *UserService) ParseToken(ctx context.Context, token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := []byte(s.hmacSecret)
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("Invalid token or claims")
	}

	return claims, nil
}
