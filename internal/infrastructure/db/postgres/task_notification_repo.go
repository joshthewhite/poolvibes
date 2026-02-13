package postgres

import (
	"context"
	"database/sql"
	"fmt"

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
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO task_notifications (id, task_id, user_id, type, due_date, sent_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (task_id, type, due_date) DO NOTHING`,
		notif.ID, notif.TaskID, notif.UserID,
		notif.Type, notif.DueDate, notif.SentAt)
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
	_, err := r.db.ExecContext(ctx, `DELETE FROM task_notifications WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting task notification: %w", err)
	}
	return nil
}
