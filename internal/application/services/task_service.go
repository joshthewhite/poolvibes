package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/josh/poolio/internal/application/command"
	"github.com/josh/poolio/internal/domain/entities"
	"github.com/josh/poolio/internal/domain/repositories"
	"github.com/josh/poolio/internal/domain/valueobjects"
)

type TaskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) List(ctx context.Context) ([]entities.Task, error) {
	return s.repo.FindAll(ctx)
}

func (s *TaskService) Get(ctx context.Context, id string) (*entities.Task, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.FindByID(ctx, uid)
}

func (s *TaskService) Create(ctx context.Context, cmd command.CreateTask) (*entities.Task, error) {
	rec, err := valueobjects.NewRecurrence(valueobjects.Frequency(cmd.RecurrenceFrequency), cmd.RecurrenceInterval)
	if err != nil {
		return nil, fmt.Errorf("recurrence: %w", err)
	}
	task := entities.NewTask(cmd.Name, cmd.Description, rec, cmd.DueDate)
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Update(ctx context.Context, cmd command.UpdateTask) (*entities.Task, error) {
	uid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	task, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	rec, err := valueobjects.NewRecurrence(valueobjects.Frequency(cmd.RecurrenceFrequency), cmd.RecurrenceInterval)
	if err != nil {
		return nil, fmt.Errorf("recurrence: %w", err)
	}
	task.Name = cmd.Name
	task.Description = cmd.Description
	task.Recurrence = rec
	task.DueDate = cmd.DueDate
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Complete(ctx context.Context, id string) (*entities.Task, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	task, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	next := task.Complete()
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("updating completed task: %w", err)
	}
	if err := s.repo.Create(ctx, next); err != nil {
		return nil, fmt.Errorf("creating next task: %w", err)
	}
	return next, nil
}

func (s *TaskService) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.Delete(ctx, uid)
}
