package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EquipmentCategory string

const (
	CategoryPump        EquipmentCategory = "pump"
	CategoryFilter      EquipmentCategory = "filter"
	CategoryHeater      EquipmentCategory = "heater"
	CategoryChlorinator EquipmentCategory = "chlorinator"
	CategoryCleaner     EquipmentCategory = "cleaner"
	CategoryOther       EquipmentCategory = "other"
)

type Equipment struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Name           string
	Category       EquipmentCategory
	Manufacturer   string
	Model          string
	SerialNumber   string
	InstallDate    *time.Time
	WarrantyExpiry *time.Time
	ServiceRecords []ServiceRecord
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewEquipment(userID uuid.UUID, name string, category EquipmentCategory, manufacturer, model, serialNumber string, installDate, warrantyExpiry *time.Time) *Equipment {
	now := time.Now()
	return &Equipment{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         userID,
		Name:           name,
		Category:       category,
		Manufacturer:   manufacturer,
		Model:          model,
		SerialNumber:   serialNumber,
		InstallDate:    installDate,
		WarrantyExpiry: warrantyExpiry,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (e *Equipment) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch e.Category {
	case CategoryPump, CategoryFilter, CategoryHeater, CategoryChlorinator, CategoryCleaner, CategoryOther:
	default:
		return fmt.Errorf("invalid category: %s", e.Category)
	}
	return nil
}

func (e *Equipment) IsWarrantyActive() bool {
	if e.WarrantyExpiry == nil {
		return false
	}
	return time.Now().Before(*e.WarrantyExpiry)
}
