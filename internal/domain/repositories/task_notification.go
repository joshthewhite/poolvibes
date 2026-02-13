package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type TaskNotificationRepository interface {
	// Claim atomically inserts a notification record. Returns true if this
	// caller claimed it, false if another instance already did (UNIQUE conflict).
	Claim(ctx context.Context, notif *entities.TaskNotification) (bool, error)
	// Delete removes a notification record (used to release a claim on send failure).
	Delete(ctx context.Context, id uuid.UUID) error
}
