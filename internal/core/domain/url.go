package domain

import "time"

type URL struct {
	ID          string
	Destination string
	UserID      int64
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
