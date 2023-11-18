package entity

import (
	"time"

	uuid "github.com/google/uuid"
)

type User struct {
	UUID         uuid.UUID `bson:"uuid"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Password     string    `bson:"password"`
	RegisteredAt time.Time `bson:"registered_at"`
}
