package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type ChemicalRepository interface {
	FindAll(ctx context.Context) ([]entities.Chemical, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Chemical, error)
	Create(ctx context.Context, chemical *entities.Chemical) error
	Update(ctx context.Context, chemical *entities.Chemical) error
	Delete(ctx context.Context, id uuid.UUID) error
}
