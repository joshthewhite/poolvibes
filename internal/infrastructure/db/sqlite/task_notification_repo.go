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

func (r *TaskNotificationRepo) Claim(ctx context.Context, notif *entities.TaskNotification) (bool, error) {
	var taskID string
	if notif.TaskID != uuid.Nil {
		taskID = notif.TaskID.String()
	}
	res, err := r.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO task_notifications (id, task_id, user_id, type, due_date, sent_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		notif.ID.String(), taskID, notif.UserID.String(),
		notif.Type, notif.DueDate.Format("2006-01-02"), notif.SentAt.Format(time.RFC3339))
	if err != nil {
		return false, fmt.Errorf("claiming task notification: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("checking rows affected: %w", err)
	}
	return rows > 0, nil
}

func (r *TaskNotificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM task_notifications WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("deleting task notification: %w", err)
	}
	return nil
}
