package entities

import (
	"time"

	"github.com/google/uuid"
)

type PushSubscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Endpoint  string
	P256dh    string
	Auth      string
	CreatedAt time.Time
}

func NewPushSubscription(userID uuid.UUID, endpoint, p256dh, auth string) *PushSubscription {
	return &PushSubscription{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    userID,
		Endpoint:  endpoint,
		P256dh:    p256dh,
		Auth:      auth,
		CreatedAt: time.Now(),
	}
}
