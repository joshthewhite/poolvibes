package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type ChemistryService struct {
	repo repositories.ChemistryLogRepository
}

func NewChemistryService(repo repositories.ChemistryLogRepository) *ChemistryService {
	return &ChemistryService{repo: repo}
}

func (s *ChemistryService) List(ctx context.Context) ([]entities.ChemistryLog, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAll(ctx, userID)
}

func (s *ChemistryService) ListPaged(ctx context.Context, query repositories.ChemistryLogQuery) (*repositories.PagedResult[entities.ChemistryLog], error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.FindPaged(ctx, userID, query)
}

func (s *ChemistryService) Get(ctx context.Context, id string) (*entities.ChemistryLog, error) {
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

func (s *ChemistryService) Create(ctx context.Context, cmd command.CreateChemistryLog) (*entities.ChemistryLog, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	log := entities.NewChemistryLog(userID, cmd.PH, cmd.FreeChlorine, cmd.CombinedChlorine, cmd.TotalAlkalinity, cmd.CYA, cmd.CalciumHardness, cmd.Temperature, cmd.Notes, cmd.TestedAt)
	if err := log.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.repo.Create(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *ChemistryService) Update(ctx context.Context, cmd command.UpdateChemistryLog) (*entities.ChemistryLog, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	log, err := s.repo.FindByID(ctx, userID, uid)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, fmt.Errorf("chemistry log not found")
	}
	log.PH = cmd.PH
	log.FreeChlorine = cmd.FreeChlorine
	log.CombinedChlorine = cmd.CombinedChlorine
	log.TotalAlkalinity = cmd.TotalAlkalinity
	log.CYA = cmd.CYA
	log.CalciumHardness = cmd.CalciumHardness
	log.Temperature = cmd.Temperature
	log.Notes = cmd.Notes
	log.TestedAt = cmd.TestedAt
	if err := log.Validate(); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.repo.Update(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *ChemistryService) Delete(ctx context.Context, id string) error {
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
