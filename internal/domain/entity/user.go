package entity

import (
	"time"
)

type User struct {
	ID           string
	Name         string
	Email        string
	Password     string
	RegisteredAt time.Time
}
