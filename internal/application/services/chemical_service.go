package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

type ChemicalService struct {
	repo repositories.ChemicalRepository
}

func NewChemicalService(repo repositories.ChemicalRepository) *ChemicalService {
	return &ChemicalService{repo: repo}
}

func (s *ChemicalService) List(ctx context.Context) ([]entities.Chemical, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAll(ctx, userID)
}

func (s *ChemicalService) Get(ctx context.Context, id string) (*entities.Chemical, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.FindByID(ctx, userID, uid)
}

func (s *ChemicalService) Create(ctx context.Context, cmd command.CreateChemical) (*entities.Chemical, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	stock, err := valueobjects.NewQuantity(cmd.StockAmount, valueobjects.Unit(cmd.StockUnit))
	if err != nil {
		return nil, fmt.Errorf("stock: %w", err)
	}
	chem := entities.NewChemical(userID, cmd.Name, entities.ChemicalType(cmd.Type), stock, cmd.AlertThreshold)
	if err := chem.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.repo.Create(ctx, chem); err != nil {
		return nil, err
	}
	return chem, nil
}

func (s *ChemicalService) Update(ctx context.Context, cmd command.UpdateChemical) (*entities.Chemical, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	chem, err := s.repo.FindByID(ctx, userID, uid)
	if err != nil {
		return nil, err
	}
	if chem == nil {
		return nil, fmt.Errorf("chemical not found")
	}
	stock, err := valueobjects.NewQuantity(cmd.StockAmount, valueobjects.Unit(cmd.StockUnit))
	if err != nil {
		return nil, fmt.Errorf("stock: %w", err)
	}
	chem.Name = cmd.Name
	chem.Type = entities.ChemicalType(cmd.Type)
	chem.Stock = stock
	chem.AlertThreshold = cmd.AlertThreshold
	if err := chem.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.repo.Update(ctx, chem); err != nil {
		return nil, err
	}
	return chem, nil
}

func (s *ChemicalService) AdjustStock(ctx context.Context, cmd command.AdjustChemicalStock) (*entities.Chemical, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	chem, err := s.repo.FindByID(ctx, userID, uid)
	if err != nil {
		return nil, err
	}
	if chem == nil {
		return nil, fmt.Errorf("chemical not found")
	}
	if err := chem.AdjustStock(cmd.Delta); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, chem); err != nil {
		return nil, err
	}
	return chem, nil
}

func (s *ChemicalService) Delete(ctx context.Context, id string) error {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return err
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.Delete(ctx, userID, uid)
}
