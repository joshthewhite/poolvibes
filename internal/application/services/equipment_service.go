package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type EquipmentService struct {
	eqRepo repositories.EquipmentRepository
	srRepo repositories.ServiceRecordRepository
}

func NewEquipmentService(eqRepo repositories.EquipmentRepository, srRepo repositories.ServiceRecordRepository) *EquipmentService {
	return &EquipmentService{eqRepo: eqRepo, srRepo: srRepo}
}

func (s *EquipmentService) List(ctx context.Context) ([]entities.Equipment, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	items, err := s.eqRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range items {
		records, err := s.srRepo.FindByEquipmentID(ctx, userID, items[i].ID)
		if err != nil {
			return nil, fmt.Errorf("loading service records: %w", err)
		}
		items[i].ServiceRecords = records
	}
	return items, nil
}

func (s *EquipmentService) Get(ctx context.Context, id string) (*entities.Equipment, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	eq, err := s.eqRepo.FindByID(ctx, userID, uid)
	if err != nil {
		return nil, err
	}
	if eq == nil {
		return nil, nil
	}
	records, err := s.srRepo.FindByEquipmentID(ctx, userID, eq.ID)
	if err != nil {
		return nil, fmt.Errorf("loading service records: %w", err)
	}
	eq.ServiceRecords = records
	return eq, nil
}

func (s *EquipmentService) Create(ctx context.Context, cmd command.CreateEquipment) (*entities.Equipment, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	eq := entities.NewEquipment(userID, cmd.Name, entities.EquipmentCategory(cmd.Category), cmd.Manufacturer, cmd.Model, cmd.SerialNumber, cmd.InstallDate, cmd.WarrantyExpiry)
	if err := eq.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.eqRepo.Create(ctx, eq); err != nil {
		return nil, err
	}
	return eq, nil
}

func (s *EquipmentService) Update(ctx context.Context, cmd command.UpdateEquipment) (*entities.Equipment, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	eq, err := s.eqRepo.FindByID(ctx, userID, uid)
	if err != nil {
		return nil, err
	}
	if eq == nil {
		return nil, fmt.Errorf("equipment not found")
	}
	eq.Name = cmd.Name
	eq.Category = entities.EquipmentCategory(cmd.Category)
	eq.Manufacturer = cmd.Manufacturer
	eq.Model = cmd.Model
	eq.SerialNumber = cmd.SerialNumber
	eq.InstallDate = cmd.InstallDate
	eq.WarrantyExpiry = cmd.WarrantyExpiry
	if err := eq.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.eqRepo.Update(ctx, eq); err != nil {
		return nil, err
	}
	return eq, nil
}

func (s *EquipmentService) Delete(ctx context.Context, id string) error {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	return s.eqRepo.Delete(ctx, userID, uid)
}

func (s *EquipmentService) AddServiceRecord(ctx context.Context, cmd command.CreateServiceRecord) (*entities.ServiceRecord, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	eqID, err := uuid.Parse(cmd.EquipmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid equipment ID: %w", err)
	}
	sr := entities.NewServiceRecord(userID, eqID, cmd.ServiceDate, cmd.Description, cmd.Cost, cmd.Technician)
	if err := sr.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.srRepo.Create(ctx, sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *EquipmentService) DeleteServiceRecord(ctx context.Context, id string) error {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	return s.srRepo.Delete(ctx, userID, uid)
}
