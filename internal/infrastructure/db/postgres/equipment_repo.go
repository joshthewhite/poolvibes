package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type EquipmentRepo struct {
	db *sql.DB
}

func NewEquipmentRepo(db *sql.DB) *EquipmentRepo {
	return &EquipmentRepo{db: db}
}

func (r *EquipmentRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Equipment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, category, manufacturer, model, serial_number, install_date, warranty_expiry, created_at, updated_at FROM equipment WHERE user_id = $1 ORDER BY name ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying equipment: %w", err)
	}
	defer rows.Close()

	var items []entities.Equipment
	for rows.Next() {
		e, err := scanEquipment(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *e)
	}
	return items, rows.Err()
}

func (r *EquipmentRepo) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.Equipment, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, category, manufacturer, model, serial_number, install_date, warranty_expiry, created_at, updated_at FROM equipment WHERE id = $1 AND user_id = $2`, id, userID)
	e, err := scanEquipmentRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying equipment: %w", err)
	}
	return e, nil
}

func (r *EquipmentRepo) Create(ctx context.Context, e *entities.Equipment) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO equipment (id, user_id, name, category, manufacturer, model, serial_number, install_date, warranty_expiry, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		e.ID, e.UserID, e.Name, string(e.Category), e.Manufacturer, e.Model, e.SerialNumber, e.InstallDate, e.WarrantyExpiry, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting equipment: %w", err)
	}
	return nil
}

func (r *EquipmentRepo) Update(ctx context.Context, e *entities.Equipment) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE equipment SET name = $1, category = $2, manufacturer = $3, model = $4, serial_number = $5, install_date = $6, warranty_expiry = $7, updated_at = $8 WHERE id = $9 AND user_id = $10`,
		e.Name, string(e.Category), e.Manufacturer, e.Model, e.SerialNumber, e.InstallDate, e.WarrantyExpiry, e.UpdatedAt, e.ID, e.UserID)
	if err != nil {
		return fmt.Errorf("updating equipment: %w", err)
	}
	return nil
}

func (r *EquipmentRepo) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM equipment WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("deleting equipment: %w", err)
	}
	return nil
}

func scanEquipmentFromRow(s scanner) (*entities.Equipment, error) {
	var e entities.Equipment
	var category string
	if err := s.Scan(&e.ID, &e.UserID, &e.Name, &category, &e.Manufacturer, &e.Model, &e.SerialNumber, &e.InstallDate, &e.WarrantyExpiry, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scanning equipment: %w", err)
	}
	e.Category = entities.EquipmentCategory(category)
	return &e, nil
}

func scanEquipment(rows *sql.Rows) (*entities.Equipment, error) {
	return scanEquipmentFromRow(rows)
}

func scanEquipmentRow(row *sql.Row) (*entities.Equipment, error) {
	return scanEquipmentFromRow(row)
}
