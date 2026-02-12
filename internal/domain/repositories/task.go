package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/josh/poolio/internal/domain/entities"
)

type TaskRepository interface {
	FindAll(ctx context.Context) ([]entities.Task, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Task, error)
	Create(ctx context.Context, task *entities.Task) error
	Update(ctx context.Context, task *entities.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
