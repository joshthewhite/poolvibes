package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type SessionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Session, error)
	Create(ctx context.Context, session *entities.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
