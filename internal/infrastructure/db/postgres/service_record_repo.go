package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type ServiceRecordRepo struct {
	db *sql.DB
}

func NewServiceRecordRepo(db *sql.DB) *ServiceRecordRepo {
	return &ServiceRecordRepo{db: db}
}

func (r *ServiceRecordRepo) FindByEquipmentID(ctx context.Context, userID uuid.UUID, equipmentID uuid.UUID) ([]entities.ServiceRecord, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, equipment_id, service_date, description, cost, technician, created_at, updated_at FROM service_records WHERE equipment_id = $1 AND user_id = $2 ORDER BY service_date DESC`, equipmentID, userID)
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

func (r *ServiceRecordRepo) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entities.ServiceRecord, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, equipment_id, service_date, description, cost, technician, created_at, updated_at FROM service_records WHERE id = $1 AND user_id = $2`, id, userID)
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
	_, err := r.db.ExecContext(ctx, `INSERT INTO service_records (id, user_id, equipment_id, service_date, description, cost, technician, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		sr.ID, sr.UserID, sr.EquipmentID, sr.ServiceDate, sr.Description, sr.Cost, sr.Technician, sr.CreatedAt, sr.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting service record: %w", err)
	}
	return nil
}

func (r *ServiceRecordRepo) Update(ctx context.Context, sr *entities.ServiceRecord) error {
	sr.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE service_records SET service_date = $1, description = $2, cost = $3, technician = $4, updated_at = $5 WHERE id = $6 AND user_id = $7`,
		sr.ServiceDate, sr.Description, sr.Cost, sr.Technician, sr.UpdatedAt, sr.ID, sr.UserID)
	if err != nil {
		return fmt.Errorf("updating service record: %w", err)
	}
	return nil
}

func (r *ServiceRecordRepo) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM service_records WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("deleting service record: %w", err)
	}
	return nil
}

func scanServiceRecordFromRow(s scanner) (*entities.ServiceRecord, error) {
	var sr entities.ServiceRecord
	if err := s.Scan(&sr.ID, &sr.UserID, &sr.EquipmentID, &sr.ServiceDate, &sr.Description, &sr.Cost, &sr.Technician, &sr.CreatedAt, &sr.UpdatedAt); err != nil {
		return nil, fmt.Errorf("scanning service record: %w", err)
	}
	return &sr, nil
}

func scanServiceRecord(rows *sql.Rows) (*entities.ServiceRecord, error) {
	return scanServiceRecordFromRow(rows)
}

func scanServiceRecordRow(row *sql.Row) (*entities.ServiceRecord, error) {
	return scanServiceRecordFromRow(row)
}
