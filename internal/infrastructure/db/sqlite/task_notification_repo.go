package sqlite

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
		WHERE task_id = ? AND type = ? AND due_date = ?`,
		taskID.String(), notifType, dueDate.Format("2006-01-02")).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking task notification: %w", err)
	}
	return count > 0, nil
}

func (r *TaskNotificationRepo) Create(ctx context.Context, notif *entities.TaskNotification) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO task_notifications (id, task_id, user_id, type, due_date, sent_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		notif.ID.String(), notif.TaskID.String(), notif.UserID.String(),
		notif.Type, notif.DueDate.Format("2006-01-02"), notif.SentAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting task notification: %w", err)
	}
	return nil
}
