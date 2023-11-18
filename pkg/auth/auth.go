package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTToken interface {
	NewJWTToken(id string) (string, error)
}

type Token struct {
	jwt jwt.Token

	secret     string
	expiration time.Time
}

func NewToken(jwt jwt.Token, secret string, expiration time.Time) *Token {
	return &Token{}
}

func (t *Token) NewJWTToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject":    id,
		"issued_at":  time.Now().Unix(),
		"expires_at": t.expiration.Unix(),
	})

	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
