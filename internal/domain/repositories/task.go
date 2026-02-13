package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type TaskRepository interface {
	FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Task, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.Task, error)
	Create(ctx context.Context, task *entities.Task) error
	Update(ctx context.Context, task *entities.Task) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}
