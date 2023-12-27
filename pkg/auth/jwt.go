package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type jwtAuthenticator struct {
	signKey string
	logger  *logrus.Logger
}

func NewAuth(logger *logrus.Logger) Authenticator {
	return &jwtAuthenticator{
		logger: logger,
	}
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
		s.logger.WithError(err).Error("failed to generate refresh token")
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(mySigningKey)
	if err != nil {
		s.logger.WithError(err).Error("failed to sign access token")
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
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
		s.logger.WithError(err).Error("failed to sign refresh token")
		return "", err
	}

	return signedRefreshToken, nil
}

func (s *jwtAuthenticator) ParseToken(accessToken string) (*ParseTokenClaimsOutput, error) {
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to parse jwt token")
		return nil, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	if !token.Valid {
		s.logger.Error("token is not valid")
		return nil, fmt.Errorf("token is not valid")
	}

	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role == nil {
		s.logger.Error("token is not valid: missing role")
		return nil, fmt.Errorf("token is not valid")
	}
	sub := claims["sub"]
	if sub == nil {
		s.logger.Error("token is not valid: missing subject")
		return nil, fmt.Errorf("token is not valid")
	}

	return &ParseTokenClaimsOutput{Sub: fmt.Sprint(sub), Role: fmt.Sprint(role)}, nil
}
