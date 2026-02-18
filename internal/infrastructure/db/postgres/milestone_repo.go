package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type MilestoneRepo struct {
	db *sql.DB
}

func NewMilestoneRepo(db *sql.DB) *MilestoneRepo {
	return &MilestoneRepo{db: db}
}

func (r *MilestoneRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Milestone, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, milestone, earned_at
		FROM user_milestones
		WHERE user_id = $1
		ORDER BY earned_at ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying milestones: %w", err)
	}
	defer rows.Close()

	var milestones []entities.Milestone
	for rows.Next() {
		var m entities.Milestone
		var milestone string
		if err := rows.Scan(&m.ID, &m.UserID, &milestone, &m.EarnedAt); err != nil {
			return nil, fmt.Errorf("scanning milestone: %w", err)
		}
		m.Milestone = entities.MilestoneKey(milestone)
		milestones = append(milestones, m)
	}
	return milestones, rows.Err()
}

func (r *MilestoneRepo) Create(ctx context.Context, m *entities.Milestone) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_milestones (id, user_id, milestone, earned_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, milestone) DO NOTHING`,
		m.ID, m.UserID, string(m.Milestone), m.EarnedAt)
	if err != nil {
		return fmt.Errorf("inserting milestone: %w", err)
	}
	return nil
}

func (r *MilestoneRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_milestones WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("deleting milestones: %w", err)
	}
	return nil
}
