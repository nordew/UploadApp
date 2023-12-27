package entity

import "time"

const (
	Upload = "upload"
	Delete = "delete"
)

type AuditLog struct {
	LogID      int
	UserID     string
	ActionType string
	OldData    []byte
	NewData    []byte
	Timestamp  time.Time
}
