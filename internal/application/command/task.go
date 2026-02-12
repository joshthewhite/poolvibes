package command

import "time"

type CreateTask struct {
	Name                string
	Description         string
	RecurrenceFrequency string
	RecurrenceInterval  int
	DueDate             time.Time
}

type UpdateTask struct {
	ID                  string
	Name                string
	Description         string
	RecurrenceFrequency string
	RecurrenceInterval  int
	DueDate             time.Time
}
