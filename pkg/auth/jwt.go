package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type jwtAuthenticator struct {
	signKey string
}

func NewAuth() Authenticator {
	return &jwtAuthenticator{}
}

type TokenClaims struct {
	UserId string `json:"sub"`
	Role   string `json:"role""`
	jwt.RegisteredClaims
}

func (s *jwtAuthenticator) GenerateTokens(options *GenerateTokenClaimsOptions) (string, string, error) {
	mySigningKey := []byte(s.signKey)

	claims := TokenClaims{
		UserId: options.UserId,
		Role:   options.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "upload-api",
			Subject:   "client",
			ID:        uuid.NewString(),
			Audience:  []string{"upload"},
		},
	}

	refreshToken, err := s.GenerateRefreshToken(options.UserId, options.Role)
	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *jwtAuthenticator) GenerateRefreshToken(id string, role string) (string, error) {
	mySigningKey := []byte(s.signKey)

	claims := TokenClaims{
		UserId: id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "upload-api",
			Subject:   "client",
			ID:        uuid.NewString(),
			Audience:  []string{"upload"},
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedRefreshToken, err := refreshToken.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return signedRefreshToken, nil
}

func (s *jwtAuthenticator) ParseToken(accessToken string) (*ParseTokenClaimsOutput, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	username := claims["username"]
	if username == nil {
		return nil, fmt.Errorf("token is not valid")
	}
	userId := claims["userId"]
	if userId == nil {
		return nil, fmt.Errorf("token is not valid")
	}

	return &ParseTokenClaimsOutput{UserId: fmt.Sprint(userId), Username: fmt.Sprint(username)}, nil
}
