package service

import (
	"context"
	"time"
)

type Subscription interface {
	Create(ctx context.Context, userID int, startDate time.Time, endDate time.Time) error
	Cancel(ctx context.Context, userID int) error
	Update(ctx context.Context, userID int, newEndDate time.Time) error
	IsExpired(ctx context.Context, userID int) (bool, time.Time, error)
}
