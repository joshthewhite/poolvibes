package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type ChemistryLogRepository interface {
	FindAll(ctx context.Context) ([]entities.ChemistryLog, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.ChemistryLog, error)
	Create(ctx context.Context, log *entities.ChemistryLog) error
	Update(ctx context.Context, log *entities.ChemistryLog) error
	Delete(ctx context.Context, id uuid.UUID) error
}
