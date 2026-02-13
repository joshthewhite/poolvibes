package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

type ChemicalRepo struct {
	db *sql.DB
}

func NewChemicalRepo(db *sql.DB) *ChemicalRepo {
	return &ChemicalRepo{db: db}
}

func (r *ChemicalRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Chemical, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, type,
			stock_amount, stock_unit, alert_threshold,
			last_purchased, created_at, updated_at
		FROM chemicals
		WHERE user_id = ?
		ORDER BY name ASC`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("querying chemicals: %w", err)
	}
	defer rows.Close()

	var chemicals []entities.Chemical
	for rows.Next() {
		c, err := scanChemical(rows)
		if err != nil {
			return nil, err
		}
		chemicals = append(chemicals, *c)
	}
	return chemicals, rows.Err()
}

func (r *ChemicalRepo) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.Chemical, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, type,
			stock_amount, stock_unit, alert_threshold,
			last_purchased, created_at, updated_at
		FROM chemicals
		WHERE id = ? AND user_id = ?`, id.String(), userID.String())
	c, err := scanChemicalRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying chemical: %w", err)
	}
	return c, nil
}

func (r *ChemicalRepo) Create(ctx context.Context, c *entities.Chemical) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO chemicals (id, user_id, name, type,
			stock_amount, stock_unit, alert_threshold,
			last_purchased, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID.String(), c.UserID.String(), c.Name, string(c.Type), c.Stock.Amount, string(c.Stock.Unit), c.AlertThreshold, fmtTimePtr(c.LastPurchased), c.CreatedAt.Format(time.RFC3339), c.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting chemical: %w", err)
	}
	return nil
}

func (r *ChemicalRepo) Update(ctx context.Context, c *entities.Chemical) error {
	c.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE chemicals
		SET name = ?, type = ?,
			stock_amount = ?, stock_unit = ?,
			alert_threshold = ?, last_purchased = ?,
			updated_at = ?
		WHERE id = ? AND user_id = ?`,
		c.Name, string(c.Type), c.Stock.Amount, string(c.Stock.Unit), c.AlertThreshold, fmtTimePtr(c.LastPurchased), c.UpdatedAt.Format(time.RFC3339), c.ID.String(), c.UserID.String())
	if err != nil {
		return fmt.Errorf("updating chemical: %w", err)
	}
	return nil
}

func (r *ChemicalRepo) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chemicals WHERE id = ? AND user_id = ?`, id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("deleting chemical: %w", err)
	}
	return nil
}

func scanChemicalFromRow(s scanner) (*entities.Chemical, error) {
	var c entities.Chemical
	var idStr, userIDStr, chemType, stockUnit, createdAt, updatedAt string
	var lastPurchased *string
	if err := s.Scan(&idStr, &userIDStr, &c.Name, &chemType, &c.Stock.Amount, &stockUnit, &c.AlertThreshold, &lastPurchased, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scanning chemical: %w", err)
	}
	c.ID = uuid.MustParse(idStr)
	c.UserID = uuid.MustParse(userIDStr)
	c.Type = entities.ChemicalType(chemType)
	c.Stock.Unit = valueobjects.Unit(stockUnit)
	c.LastPurchased = parseTimePtr(lastPurchased)
	c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &c, nil
}

func scanChemical(rows *sql.Rows) (*entities.Chemical, error) {
	return scanChemicalFromRow(rows)
}

func scanChemicalRow(row *sql.Row) (*entities.Chemical, error) {
	return scanChemicalFromRow(row)
}
