package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type EquipmentRepository interface {
	FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Equipment, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.Equipment, error)
	Create(ctx context.Context, equipment *entities.Equipment) error
	Update(ctx context.Context, equipment *entities.Equipment) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type ServiceRecordRepository interface {
	FindByEquipmentID(ctx context.Context, userID uuid.UUID, equipmentID uuid.UUID) ([]entities.ServiceRecord, error)
	FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.ServiceRecord, error)
	Create(ctx context.Context, record *entities.ServiceRecord) error
	Update(ctx context.Context, record *entities.ServiceRecord) error
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}
