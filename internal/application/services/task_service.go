package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

type TaskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) List(ctx context.Context) ([]entities.Task, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAll(ctx, userID)
}

func (s *TaskService) Get(ctx context.Context, id string) (*entities.Task, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.FindByID(ctx, userID, uid)
}

func (s *TaskService) Create(ctx context.Context, cmd command.CreateTask) (*entities.Task, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := valueobjects.NewRecurrence(valueobjects.Frequency(cmd.RecurrenceFrequency), cmd.RecurrenceInterval)
	if err != nil {
		return nil, fmt.Errorf("recurrence: %w", err)
	}
	task := entities.NewTask(userID, cmd.Name, cmd.Description, rec, cmd.DueDate)
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Update(ctx context.Context, cmd command.UpdateTask) (*entities.Task, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	task, err := s.repo.FindByID(ctx, userID, uid)
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
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	task, err := s.repo.FindByID(ctx, userID, uid)
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
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.Delete(ctx, userID, uid)
}
