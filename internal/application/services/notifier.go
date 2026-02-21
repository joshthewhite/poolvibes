package services

import (
	"context"

	"github.com/google/uuid"
)

type Notifier interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

type PushNotifier interface {
	SendToUser(ctx context.Context, userID uuid.UUID, subject string, body string) error
}
