package entity

import (
	"time"
)

type User struct {
	ID           int       `bson:"id"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Password     string    `bson:"password"`
	RegisteredAt time.Time `bson:"registered_at"`
}
