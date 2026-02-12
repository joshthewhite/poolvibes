package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/josh/poolio/internal/domain/entities"
	"github.com/josh/poolio/internal/domain/valueobjects"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) FindAll(ctx context.Context) ([]entities.Task, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, recurrence_frequency, recurrence_interval, due_date, status, completed_at, created_at, updated_at FROM tasks ORDER BY due_date ASC`)
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

func (r *TaskRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.Task, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, description, recurrence_frequency, recurrence_interval, due_date, status, completed_at, created_at, updated_at FROM tasks WHERE id = ?`, id.String())
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
	var completedAt *string
	if t.CompletedAt != nil {
		s := t.CompletedAt.Format(time.RFC3339)
		completedAt = &s
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO tasks (id, name, description, recurrence_frequency, recurrence_interval, due_date, status, completed_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID.String(), t.Name, t.Description, string(t.Recurrence.Frequency), t.Recurrence.Interval, t.DueDate.Format(time.RFC3339), string(t.Status), completedAt, t.CreatedAt.Format(time.RFC3339), t.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting task: %w", err)
	}
	return nil
}

func (r *TaskRepo) Update(ctx context.Context, t *entities.Task) error {
	t.UpdatedAt = time.Now()
	var completedAt *string
	if t.CompletedAt != nil {
		s := t.CompletedAt.Format(time.RFC3339)
		completedAt = &s
	}
	_, err := r.db.ExecContext(ctx, `UPDATE tasks SET name = ?, description = ?, recurrence_frequency = ?, recurrence_interval = ?, due_date = ?, status = ?, completed_at = ?, updated_at = ? WHERE id = ?`,
		t.Name, t.Description, string(t.Recurrence.Frequency), t.Recurrence.Interval, t.DueDate.Format(time.RFC3339), string(t.Status), completedAt, t.UpdatedAt.Format(time.RFC3339), t.ID.String())
	if err != nil {
		return fmt.Errorf("updating task: %w", err)
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTaskFromRow(s scanner) (*entities.Task, error) {
	var t entities.Task
	var idStr, freq, dueDate, status, createdAt, updatedAt string
	var interval int
	var completedAt *string
	if err := s.Scan(&idStr, &t.Name, &t.Description, &freq, &interval, &dueDate, &status, &completedAt, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	t.ID = uuid.MustParse(idStr)
	t.Recurrence = valueobjects.Recurrence{Frequency: valueobjects.Frequency(freq), Interval: interval}
	t.DueDate, _ = time.Parse(time.RFC3339, dueDate)
	t.Status = entities.TaskStatus(status)
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	if completedAt != nil {
		ca, _ := time.Parse(time.RFC3339, *completedAt)
		t.CompletedAt = &ca
	}
	return &t, nil
}

func scanTask(rows *sql.Rows) (*entities.Task, error) {
	return scanTaskFromRow(rows)
}

func scanTaskRow(row *sql.Row) (*entities.Task, error) {
	return scanTaskFromRow(row)
}
