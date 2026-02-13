package sqlite

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
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, category,
			manufacturer, model, serial_number,
			install_date, warranty_expiry,
			created_at, updated_at
		FROM equipment
		WHERE user_id = ?
		ORDER BY name ASC`, userID.String())
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
	row := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, name, category,
			manufacturer, model, serial_number,
			install_date, warranty_expiry,
			created_at, updated_at
		FROM equipment
		WHERE id = ? AND user_id = ?`, id.String(), userID.String())
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
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO equipment (id, user_id, name, category,
			manufacturer, model, serial_number,
			install_date, warranty_expiry,
			created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID.String(), e.UserID.String(), e.Name, string(e.Category), e.Manufacturer, e.Model, e.SerialNumber, fmtTimePtr(e.InstallDate), fmtTimePtr(e.WarrantyExpiry), e.CreatedAt.Format(time.RFC3339), e.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting equipment: %w", err)
	}
	return nil
}

func (r *EquipmentRepo) Update(ctx context.Context, e *entities.Equipment) error {
	e.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE equipment
		SET name = ?, category = ?,
			manufacturer = ?, model = ?, serial_number = ?,
			install_date = ?, warranty_expiry = ?,
			updated_at = ?
		WHERE id = ? AND user_id = ?`,
		e.Name, string(e.Category), e.Manufacturer, e.Model, e.SerialNumber, fmtTimePtr(e.InstallDate), fmtTimePtr(e.WarrantyExpiry), e.UpdatedAt.Format(time.RFC3339), e.ID.String(), e.UserID.String())
	if err != nil {
		return fmt.Errorf("updating equipment: %w", err)
	}
	return nil
}

func (r *EquipmentRepo) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM equipment WHERE id = ? AND user_id = ?`, id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("deleting equipment: %w", err)
	}
	return nil
}

func fmtTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t, _ := time.Parse(time.RFC3339, *s)
	return &t
}

func scanEquipmentFromRow(s scanner) (*entities.Equipment, error) {
	var e entities.Equipment
	var idStr, userIDStr, category, createdAt, updatedAt string
	var installDate, warrantyExpiry *string
	if err := s.Scan(&idStr, &userIDStr, &e.Name, &category, &e.Manufacturer, &e.Model, &e.SerialNumber, &installDate, &warrantyExpiry, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scanning equipment: %w", err)
	}
	e.ID = uuid.MustParse(idStr)
	e.UserID = uuid.MustParse(userIDStr)
	e.Category = entities.EquipmentCategory(category)
	e.InstallDate = parseTimePtr(installDate)
	e.WarrantyExpiry = parseTimePtr(warrantyExpiry)
	e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	e.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &e, nil
}

func scanEquipment(rows *sql.Rows) (*entities.Equipment, error) {
	return scanEquipmentFromRow(rows)
}

func scanEquipmentRow(row *sql.Row) (*entities.Equipment, error) {
	return scanEquipmentFromRow(row)
}
