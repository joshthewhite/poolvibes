package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Task, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, description, recurrence_frequency, recurrence_interval, due_date, status, completed_at, created_at, updated_at FROM tasks WHERE user_id = $1 ORDER BY due_date ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying tasks: %w", err)
	}
	defer rows.Close()

	var tasks []entities.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		t.CheckOverdue()
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepo) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.Task, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, description, recurrence_frequency, recurrence_interval, due_date, status, completed_at, created_at, updated_at FROM tasks WHERE id = $1 AND user_id = $2`, id, userID)
	t, err := scanTaskRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying task: %w", err)
	}
	t.CheckOverdue()
	return t, nil
}

func (r *TaskRepo) Create(ctx context.Context, t *entities.Task) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO tasks (id, user_id, name, description, recurrence_frequency, recurrence_interval, due_date, status, completed_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		t.ID, t.UserID, t.Name, t.Description, string(t.Recurrence.Frequency), t.Recurrence.Interval, t.DueDate, string(t.Status), t.CompletedAt, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting task: %w", err)
	}
	return nil
}

func (r *TaskRepo) Update(ctx context.Context, t *entities.Task) error {
	t.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE tasks SET name = $1, description = $2, recurrence_frequency = $3, recurrence_interval = $4, due_date = $5, status = $6, completed_at = $7, updated_at = $8 WHERE id = $9 AND user_id = $10`,
		t.Name, t.Description, string(t.Recurrence.Frequency), t.Recurrence.Interval, t.DueDate, string(t.Status), t.CompletedAt, t.UpdatedAt, t.ID, t.UserID)
	if err != nil {
		return fmt.Errorf("updating task: %w", err)
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	return nil
}

func scanTaskFromRow(s scanner) (*entities.Task, error) {
	var t entities.Task
	var freq, status string
	if err := s.Scan(&t.ID, &t.UserID, &t.Name, &t.Description, &freq, &t.Recurrence.Interval, &t.DueDate, &status, &t.CompletedAt, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	t.Recurrence.Frequency = valueobjects.Frequency(freq)
	t.Status = entities.TaskStatus(status)
	return &t, nil
}

func scanTask(rows *sql.Rows) (*entities.Task, error) {
	return scanTaskFromRow(rows)
}

func scanTaskRow(row *sql.Row) (*entities.Task, error) {
	return scanTaskFromRow(row)
}
