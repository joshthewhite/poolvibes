package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/josh/poolio/internal/domain/valueobjects"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusOverdue   TaskStatus = "overdue"
)

type Task struct {
	ID          uuid.UUID
	Name        string
	Description string
	Recurrence  valueobjects.Recurrence
	DueDate     time.Time
	Status      TaskStatus
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(name, description string, recurrence valueobjects.Recurrence, dueDate time.Time) *Task {
	now := time.Now()
	return &Task{
		ID:          uuid.Must(uuid.NewV7()),
		Name:        name,
		Description: description,
		Recurrence:  recurrence,
		DueDate:     dueDate,
		Status:      TaskStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	if t.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}
	return nil
}

func (t *Task) Complete() *Task {
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now

	next := NewTask(t.Name, t.Description, t.Recurrence, t.Recurrence.NextDueDate(now))
	return next
}

func (t *Task) CheckOverdue() {
	if t.Status == TaskStatusPending && time.Now().After(t.DueDate) {
		t.Status = TaskStatusOverdue
	}
}
