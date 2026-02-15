package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type ChemistryLogRepo struct {
	db *sql.DB
}

func NewChemistryLogRepo(db *sql.DB) *ChemistryLogRepo {
	return &ChemistryLogRepo{db: db}
}

func (r *ChemistryLogRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.ChemistryLog, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, ph, free_chlorine, combined_chlorine,
			total_alkalinity, cya, calcium_hardness,
			temperature, notes, tested_at,
			created_at, updated_at
		FROM chemistry_logs
		WHERE user_id = $1
		ORDER BY tested_at DESC`, userID)
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

// allowedSortColumns prevents SQL injection in ORDER BY clauses.
var allowedSortColumns = map[string]string{
	"tested_at":        "tested_at",
	"ph":               "ph",
	"free_chlorine":    "free_chlorine",
	"total_alkalinity": "total_alkalinity",
	"cya":              "cya",
}

const outOfRangeWhere = `(ph < 7.2 OR ph > 7.6 OR free_chlorine < 1.0 OR free_chlorine > 3.0 OR combined_chlorine > 0.5 OR total_alkalinity < 80 OR total_alkalinity > 120 OR cya < 30 OR cya > 50 OR calcium_hardness < 200 OR calcium_hardness > 400)`

func (r *ChemistryLogRepo) FindPaged(ctx context.Context, userID uuid.UUID, query repositories.ChemistryLogQuery) (*repositories.PagedResult[entities.ChemistryLog], error) {
	query.Defaults()

	var where []string
	var args []any
	paramN := 1

	where = append(where, fmt.Sprintf("user_id = $%d", paramN))
	args = append(args, userID)
	paramN++

	if query.DateFrom != nil {
		where = append(where, fmt.Sprintf("tested_at >= $%d", paramN))
		args = append(args, *query.DateFrom)
		paramN++
	}
	if query.DateTo != nil {
		where = append(where, fmt.Sprintf("tested_at <= $%d", paramN))
		args = append(args, *query.DateTo)
		paramN++
	}
	if query.OutOfRange {
		where = append(where, outOfRangeWhere)
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	// Count total
	var total int
	countSQL := "SELECT COUNT(*) FROM chemistry_logs " + whereClause
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("counting chemistry logs: %w", err)
	}

	totalPages := (total + query.PageSize - 1) / query.PageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// Sort column
	sortCol := "tested_at"
	if col, ok := allowedSortColumns[query.SortBy]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if query.SortDir == repositories.SortAsc {
		sortDir = "ASC"
	}

	dataSQL := fmt.Sprintf(`
		SELECT id, user_id, ph, free_chlorine, combined_chlorine,
			total_alkalinity, cya, calcium_hardness,
			temperature, notes, tested_at,
			created_at, updated_at
		FROM chemistry_logs
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`, whereClause, sortCol, sortDir, paramN, paramN+1)

	dataArgs := append(args, query.PageSize, query.Offset())
	rows, err := r.db.QueryContext(ctx, dataSQL, dataArgs...)
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating chemistry logs: %w", err)
	}

	return &repositories.PagedResult[entities.ChemistryLog]{
		Items:      logs,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (r *ChemistryLogRepo) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.ChemistryLog, error) {
	var l entities.ChemistryLog
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, ph, free_chlorine, combined_chlorine,
			total_alkalinity, cya, calcium_hardness,
			temperature, notes, tested_at,
			created_at, updated_at
		FROM chemistry_logs
		WHERE id = $1 AND user_id = $2`, id, userID).
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
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO chemistry_logs (id, user_id, ph, free_chlorine,
			combined_chlorine, total_alkalinity, cya, calcium_hardness,
			temperature, notes, tested_at,
			created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		l.ID, l.UserID, l.PH, l.FreeChlorine, l.CombinedChlorine, l.TotalAlkalinity, l.CYA, l.CalciumHardness, l.Temperature, l.Notes, l.TestedAt, l.CreatedAt, l.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting chemistry log: %w", err)
	}
	return nil
}

func (r *ChemistryLogRepo) Update(ctx context.Context, l *entities.ChemistryLog) error {
	l.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE chemistry_logs
		SET ph = $1, free_chlorine = $2, combined_chlorine = $3,
			total_alkalinity = $4, cya = $5, calcium_hardness = $6,
			temperature = $7, notes = $8, tested_at = $9,
			updated_at = $10
		WHERE id = $11 AND user_id = $12`,
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
