package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr string
	}{
		{
			name:    "missing name",
			task:    Task{DueDate: time.Now()},
			wantErr: "name is required",
		},
		{
			name:    "zero due date",
			task:    Task{Name: "Test task"},
			wantErr: "due date is required",
		},
		{
			name: "valid",
			task: Task{Name: "Test task", DueDate: time.Now()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantErr {
					t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTask_CheckOverdue(t *testing.T) {
	tests := []struct {
		name       string
		status     TaskStatus
		dueDate    time.Time
		wantStatus TaskStatus
	}{
		{
			name:       "pending past due becomes overdue",
			status:     TaskStatusPending,
			dueDate:    time.Now().Add(-24 * time.Hour),
			wantStatus: TaskStatusOverdue,
		},
		{
			name:       "pending future due stays pending",
			status:     TaskStatusPending,
			dueDate:    time.Now().Add(24 * time.Hour),
			wantStatus: TaskStatusPending,
		},
		{
			name:       "completed past due stays completed",
			status:     TaskStatusCompleted,
			dueDate:    time.Now().Add(-24 * time.Hour),
			wantStatus: TaskStatusCompleted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Status: tt.status, DueDate: tt.dueDate}
			task.CheckOverdue()
			if task.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", task.Status, tt.wantStatus)
			}
		})
	}
}

func TestTask_Complete(t *testing.T) {
	userID := uuid.New()
	rec, _ := valueobjects.NewRecurrence(valueobjects.FrequencyDaily, 1)
	dueDate := time.Date(2025, 3, 1, 12, 0, 0, 0, time.UTC)

	task := &Task{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "Skim pool",
		Description: "Skim debris",
		Recurrence:  rec,
		DueDate:     dueDate,
		Status:      TaskStatusPending,
	}

	next := task.Complete()

	if task.Status != TaskStatusCompleted {
		t.Errorf("original Status = %v, want %v", task.Status, TaskStatusCompleted)
	}
	if task.CompletedAt == nil {
		t.Fatal("original CompletedAt should be set")
	}
	if next.Name != task.Name {
		t.Errorf("next Name = %v, want %v", next.Name, task.Name)
	}
	if next.UserID != userID {
		t.Errorf("next UserID = %v, want %v", next.UserID, userID)
	}
	if next.Status != TaskStatusPending {
		t.Errorf("next Status = %v, want %v", next.Status, TaskStatusPending)
	}
	expectedDue := time.Date(2025, 3, 2, 12, 0, 0, 0, time.UTC)
	if !next.DueDate.Equal(expectedDue) {
		t.Errorf("next DueDate = %v, want %v", next.DueDate, expectedDue)
	}
}

func TestTask_Complete_WeeklyAndMonthly(t *testing.T) {
	userID := uuid.New()
	dueDate := time.Date(2025, 3, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		freq    valueobjects.Frequency
		intv    int
		wantDue time.Time
	}{
		{"weekly", valueobjects.FrequencyWeekly, 1, time.Date(2025, 3, 8, 12, 0, 0, 0, time.UTC)},
		{"biweekly", valueobjects.FrequencyWeekly, 2, time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)},
		{"monthly", valueobjects.FrequencyMonthly, 1, time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec, _ := valueobjects.NewRecurrence(tt.freq, tt.intv)
			task := &Task{
				UserID:     userID,
				Name:       "Task",
				Recurrence: rec,
				DueDate:    dueDate,
				Status:     TaskStatusPending,
			}
			next := task.Complete()
			if !next.DueDate.Equal(tt.wantDue) {
				t.Errorf("next DueDate = %v, want %v", next.DueDate, tt.wantDue)
			}
		})
	}
}
