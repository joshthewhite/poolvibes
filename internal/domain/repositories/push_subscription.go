package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type PushSubscriptionRepository interface {
	Save(ctx context.Context, sub *entities.PushSubscription) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.PushSubscription, error)
	DeleteByEndpoint(ctx context.Context, userID uuid.UUID, endpoint string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
