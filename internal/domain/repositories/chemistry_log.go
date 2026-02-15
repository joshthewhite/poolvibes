package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type ChemistryLogRepository interface {
	FindAll(ctx context.Context, userID uuid.UUID) ([]entities.ChemistryLog, error)
	FindPaged(ctx context.Context, userID uuid.UUID, query ChemistryLogQuery) (*PagedResult[entities.ChemistryLog], error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.ChemistryLog, error)
	Create(ctx context.Context, log *entities.ChemistryLog) error
	Update(ctx context.Context, log *entities.ChemistryLog) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}
