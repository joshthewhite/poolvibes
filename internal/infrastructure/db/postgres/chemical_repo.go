package postgres

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
		WHERE user_id = $1
		ORDER BY name ASC`, userID)
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
		WHERE id = $1 AND user_id = $2`, id, userID)
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		c.ID, c.UserID, c.Name, string(c.Type), c.Stock.Amount, string(c.Stock.Unit), c.AlertThreshold, c.LastPurchased, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting chemical: %w", err)
	}
	return nil
}

func (r *ChemicalRepo) Update(ctx context.Context, c *entities.Chemical) error {
	c.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE chemicals
		SET name = $1, type = $2,
			stock_amount = $3, stock_unit = $4,
			alert_threshold = $5, last_purchased = $6,
			updated_at = $7
		WHERE id = $8 AND user_id = $9`,
		c.Name, string(c.Type), c.Stock.Amount, string(c.Stock.Unit), c.AlertThreshold, c.LastPurchased, c.UpdatedAt, c.ID, c.UserID)
	if err != nil {
		return fmt.Errorf("updating chemical: %w", err)
	}
	return nil
}

func (r *ChemicalRepo) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chemicals WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("deleting chemical: %w", err)
	}
	return nil
}

func scanChemicalFromRow(s scanner) (*entities.Chemical, error) {
	var c entities.Chemical
	var chemType, stockUnit string
	if err := s.Scan(&c.ID, &c.UserID, &c.Name, &chemType, &c.Stock.Amount, &stockUnit, &c.AlertThreshold, &c.LastPurchased, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scanning chemical: %w", err)
	}
	c.Type = entities.ChemicalType(chemType)
	c.Stock.Unit = valueobjects.Unit(stockUnit)
	return &c, nil
}

func scanChemical(rows *sql.Rows) (*entities.Chemical, error) {
	return scanChemicalFromRow(rows)
}

func scanChemicalRow(row *sql.Row) (*entities.Chemical, error) {
	return scanChemicalFromRow(row)
}
