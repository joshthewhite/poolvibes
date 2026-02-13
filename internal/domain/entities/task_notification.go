package entities

import (
	"time"

	"github.com/google/uuid"
)

type TaskNotification struct {
	ID      uuid.UUID
	TaskID  uuid.UUID
	UserID  uuid.UUID
	Type    string // "email" or "sms"
	DueDate time.Time
	SentAt  time.Time
}

func NewTaskNotification(taskID, userID uuid.UUID, notifType string, dueDate time.Time) *TaskNotification {
	return &TaskNotification{
		ID:      uuid.Must(uuid.NewV7()),
		TaskID:  taskID,
		UserID:  userID,
		Type:    notifType,
		DueDate: dueDate,
		SentAt:  time.Now(),
	}
}
