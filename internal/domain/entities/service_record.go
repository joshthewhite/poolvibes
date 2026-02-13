package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ServiceRecord struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	EquipmentID uuid.UUID
	ServiceDate time.Time
	Description string
	Cost        float64
	Technician  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewServiceRecord(userID, equipmentID uuid.UUID, serviceDate time.Time, description string, cost float64, technician string) *ServiceRecord {
	now := time.Now()
	return &ServiceRecord{
		ID:          uuid.Must(uuid.NewV7()),
		UserID:      userID,
		EquipmentID: equipmentID,
		ServiceDate: serviceDate,
		Description: description,
		Cost:        cost,
		Technician:  technician,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (s *ServiceRecord) Validate() error {
	if s.EquipmentID == uuid.Nil {
		return fmt.Errorf("equipment ID is required")
	}
	if s.ServiceDate.IsZero() {
		return fmt.Errorf("service date is required")
	}
	if s.Description == "" {
		return fmt.Errorf("description is required")
	}
	if s.Cost < 0 {
		return fmt.Errorf("cost cannot be negative")
	}
	return nil
}
