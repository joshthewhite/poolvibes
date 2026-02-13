package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewSession(userID uuid.UUID, duration time.Duration) *Session {
	now := time.Now()
	return &Session{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    userID,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}
}

func (s *Session) Validate() error {
	if s.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}
	if s.ExpiresAt.IsZero() {
		return fmt.Errorf("expiry is required")
	}
	return nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
