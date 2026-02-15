package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]entities.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindExpiredDemo(ctx context.Context, now time.Time) ([]entities.User, error)
	CountDemo(ctx context.Context) (int, error)
	CountAdmins(ctx context.Context) (int, error)
}
