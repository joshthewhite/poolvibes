package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type TaskNotificationRepository interface {
	ExistsByTaskAndType(ctx context.Context, taskID uuid.UUID, notifType string, dueDate time.Time) (bool, error)
	Create(ctx context.Context, notif *entities.TaskNotification) error
}
