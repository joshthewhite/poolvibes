package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type TaskNotificationRepo struct {
	db *sql.DB
}

func NewTaskNotificationRepo(db *sql.DB) *TaskNotificationRepo {
	return &TaskNotificationRepo{db: db}
}

func (r *TaskNotificationRepo) ExistsByTaskAndType(ctx context.Context, taskID uuid.UUID, notifType string, dueDate time.Time) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM task_notifications
		WHERE task_id = $1 AND type = $2 AND due_date = $3`,
		taskID, notifType, dueDate).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking task notification: %w", err)
	}
	return count > 0, nil
}

func (r *TaskNotificationRepo) Create(ctx context.Context, notif *entities.TaskNotification) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO task_notifications (id, task_id, user_id, type, due_date, sent_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		notif.ID, notif.TaskID, notif.UserID,
		notif.Type, notif.DueDate, notif.SentAt)
	if err != nil {
		return fmt.Errorf("inserting task notification: %w", err)
	}
	return nil
}
