package hasher

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher interface {
	// Hash hashes the given password string and returns the hashed string.
	Hash(password string) (string, error)
}

type hasher struct {
	salt string
}

func NewPasswordHasher(salt string) *hasher {
	return &hasher{salt: salt}
}

func (h *hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
