package entity

import "time"

type RefreshTokenSession struct {
	UserID       string
	Role         string
	RefreshToken string
	ExpiresAt    time.Time
}
