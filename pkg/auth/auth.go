package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Token interface {
	// NewJWTToken creates a new JSON Web Token (JWT) with the given ID.
	// It returns the token string and an error if the operation fails.
	NewToken(id string) (string, error)
}

type Auth struct {
	jwt jwt.Token

	secret     string
	expiration time.Time
}

func NewAuth(jwt jwt.Token, secret string, expiration time.Time) *Auth {
	return &Auth{}
}

func (t *Auth) NewToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   id,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: t.expiration.Unix(),
	})

	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
