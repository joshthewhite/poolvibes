package postgres

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

func (r *ChemistryLogRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.ChemistryLog, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, ph, free_chlorine, combined_chlorine, total_alkalinity, cya, calcium_hardness, temperature, notes, tested_at, created_at, updated_at FROM chemistry_logs WHERE user_id = $1 ORDER BY tested_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying chemistry logs: %w", err)
	}
	defer rows.Close()

	var logs []entities.ChemistryLog
	for rows.Next() {
		var l entities.ChemistryLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.PH, &l.FreeChlorine, &l.CombinedChlorine, &l.TotalAlkalinity, &l.CYA, &l.CalciumHardness, &l.Temperature, &l.Notes, &l.TestedAt, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning chemistry log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *ChemistryLogRepo) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.ChemistryLog, error) {
	var l entities.ChemistryLog
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, ph, free_chlorine, combined_chlorine, total_alkalinity, cya, calcium_hardness, temperature, notes, tested_at, created_at, updated_at FROM chemistry_logs WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&l.ID, &l.UserID, &l.PH, &l.FreeChlorine, &l.CombinedChlorine, &l.TotalAlkalinity, &l.CYA, &l.CalciumHardness, &l.Temperature, &l.Notes, &l.TestedAt, &l.CreatedAt, &l.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying chemistry log: %w", err)
	}
	return &l, nil
}

func (r *ChemistryLogRepo) Create(ctx context.Context, l *entities.ChemistryLog) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO chemistry_logs (id, user_id, ph, free_chlorine, combined_chlorine, total_alkalinity, cya, calcium_hardness, temperature, notes, tested_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		l.ID, l.UserID, l.PH, l.FreeChlorine, l.CombinedChlorine, l.TotalAlkalinity, l.CYA, l.CalciumHardness, l.Temperature, l.Notes, l.TestedAt, l.CreatedAt, l.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting chemistry log: %w", err)
	}
	return nil
}

func (r *ChemistryLogRepo) Update(ctx context.Context, l *entities.ChemistryLog) error {
	l.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE chemistry_logs SET ph = $1, free_chlorine = $2, combined_chlorine = $3, total_alkalinity = $4, cya = $5, calcium_hardness = $6, temperature = $7, notes = $8, tested_at = $9, updated_at = $10 WHERE id = $11 AND user_id = $12`,
		l.PH, l.FreeChlorine, l.CombinedChlorine, l.TotalAlkalinity, l.CYA, l.CalciumHardness, l.Temperature, l.Notes, l.TestedAt, l.UpdatedAt, l.ID, l.UserID)
	if err != nil {
		return fmt.Errorf("updating chemistry log: %w", err)
	}
	return nil
}

func (r *ChemistryLogRepo) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chemistry_logs WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("deleting chemistry log: %w", err)
	}
	return nil
}
