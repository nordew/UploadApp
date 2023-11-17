package hasher

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *PasswordHasher {
	return &PasswordHasher{salt: salt}
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
