package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type ChemistryLogRepo struct {
	db *sql.DB
}

func NewChemistryLogRepo(db *sql.DB) *ChemistryLogRepo {
	return &ChemistryLogRepo{db: db}
}

func (r *ChemistryLogRepo) FindAll(ctx context.Context) ([]entities.ChemistryLog, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, ph, free_chlorine, combined_chlorine, total_alkalinity, cya, calcium_hardness, temperature, notes, tested_at, created_at, updated_at FROM chemistry_logs ORDER BY tested_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("querying chemistry logs: %w", err)
	}
	defer rows.Close()

	var logs []entities.ChemistryLog
	for rows.Next() {
		var l entities.ChemistryLog
		var idStr, testedAt, createdAt, updatedAt string
		if err := rows.Scan(&idStr, &l.PH, &l.FreeChlorine, &l.CombinedChlorine, &l.TotalAlkalinity, &l.CYA, &l.CalciumHardness, &l.Temperature, &l.Notes, &testedAt, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("scanning chemistry log: %w", err)
		}
		l.ID = uuid.MustParse(idStr)
		l.TestedAt, _ = time.Parse(time.RFC3339, testedAt)
		l.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		l.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *ChemistryLogRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.ChemistryLog, error) {
	var l entities.ChemistryLog
	var idStr, testedAt, createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx, `SELECT id, ph, free_chlorine, combined_chlorine, total_alkalinity, cya, calcium_hardness, temperature, notes, tested_at, created_at, updated_at FROM chemistry_logs WHERE id = ?`, id.String()).
		Scan(&idStr, &l.PH, &l.FreeChlorine, &l.CombinedChlorine, &l.TotalAlkalinity, &l.CYA, &l.CalciumHardness, &l.Temperature, &l.Notes, &testedAt, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying chemistry log: %w", err)
	}
	l.ID = uuid.MustParse(idStr)
	l.TestedAt, _ = time.Parse(time.RFC3339, testedAt)
	l.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	l.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &l, nil
}

func (r *ChemistryLogRepo) Create(ctx context.Context, l *entities.ChemistryLog) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO chemistry_logs (id, ph, free_chlorine, combined_chlorine, total_alkalinity, cya, calcium_hardness, temperature, notes, tested_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		l.ID.String(), l.PH, l.FreeChlorine, l.CombinedChlorine, l.TotalAlkalinity, l.CYA, l.CalciumHardness, l.Temperature, l.Notes, l.TestedAt.Format(time.RFC3339), l.CreatedAt.Format(time.RFC3339), l.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting chemistry log: %w", err)
	}
	return nil
}

func (r *ChemistryLogRepo) Update(ctx context.Context, l *entities.ChemistryLog) error {
	l.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE chemistry_logs SET ph = ?, free_chlorine = ?, combined_chlorine = ?, total_alkalinity = ?, cya = ?, calcium_hardness = ?, temperature = ?, notes = ?, tested_at = ?, updated_at = ? WHERE id = ?`,
		l.PH, l.FreeChlorine, l.CombinedChlorine, l.TotalAlkalinity, l.CYA, l.CalciumHardness, l.Temperature, l.Notes, l.TestedAt.Format(time.RFC3339), l.UpdatedAt.Format(time.RFC3339), l.ID.String())
	if err != nil {
		return fmt.Errorf("updating chemistry log: %w", err)
	}
	return nil
}

func (r *ChemistryLogRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chemistry_logs WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("deleting chemistry log: %w", err)
	}
	return nil
}
