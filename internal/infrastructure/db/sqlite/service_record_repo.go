package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/josh/poolio/internal/domain/entities"
)

type ServiceRecordRepo struct {
	db *sql.DB
}

func NewServiceRecordRepo(db *sql.DB) *ServiceRecordRepo {
	return &ServiceRecordRepo{db: db}
}

func (r *ServiceRecordRepo) FindByEquipmentID(ctx context.Context, equipmentID uuid.UUID) ([]entities.ServiceRecord, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, equipment_id, service_date, description, cost, technician, created_at, updated_at FROM service_records WHERE equipment_id = ? ORDER BY service_date DESC`, equipmentID.String())
	if err != nil {
		return nil, fmt.Errorf("querying service records: %w", err)
	}
	defer rows.Close()

	var records []entities.ServiceRecord
	for rows.Next() {
		sr, err := scanServiceRecord(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *sr)
	}
	return records, rows.Err()
}

func (r *ServiceRecordRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.ServiceRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, equipment_id, service_date, description, cost, technician, created_at, updated_at FROM service_records WHERE id = ?`, id.String())
	sr, err := scanServiceRecordRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying service record: %w", err)
	}
	return sr, nil
}

func (r *ServiceRecordRepo) Create(ctx context.Context, sr *entities.ServiceRecord) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO service_records (id, equipment_id, service_date, description, cost, technician, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		sr.ID.String(), sr.EquipmentID.String(), sr.ServiceDate.Format(time.RFC3339), sr.Description, sr.Cost, sr.Technician, sr.CreatedAt.Format(time.RFC3339), sr.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting service record: %w", err)
	}
	return nil
}

func (r *ServiceRecordRepo) Update(ctx context.Context, sr *entities.ServiceRecord) error {
	sr.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE service_records SET service_date = ?, description = ?, cost = ?, technician = ?, updated_at = ? WHERE id = ?`,
		sr.ServiceDate.Format(time.RFC3339), sr.Description, sr.Cost, sr.Technician, sr.UpdatedAt.Format(time.RFC3339), sr.ID.String())
	if err != nil {
		return fmt.Errorf("updating service record: %w", err)
	}
	return nil
}

func (r *ServiceRecordRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM service_records WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("deleting service record: %w", err)
	}
	return nil
}

func scanServiceRecordFromRow(s scanner) (*entities.ServiceRecord, error) {
	var sr entities.ServiceRecord
	var idStr, eqIDStr, serviceDate, createdAt, updatedAt string
	if err := s.Scan(&idStr, &eqIDStr, &serviceDate, &sr.Description, &sr.Cost, &sr.Technician, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scanning service record: %w", err)
	}
	sr.ID = uuid.MustParse(idStr)
	sr.EquipmentID = uuid.MustParse(eqIDStr)
	sr.ServiceDate, _ = time.Parse(time.RFC3339, serviceDate)
	sr.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	sr.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &sr, nil
}

func scanServiceRecord(rows *sql.Rows) (*entities.ServiceRecord, error) {
	return scanServiceRecordFromRow(rows)
}

func scanServiceRecordRow(row *sql.Row) (*entities.ServiceRecord, error) {
	return scanServiceRecordFromRow(row)
}
