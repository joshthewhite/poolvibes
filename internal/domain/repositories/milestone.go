package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type MilestoneRepository interface {
	FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Milestone, error)
	Create(ctx context.Context, milestone *entities.Milestone) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
