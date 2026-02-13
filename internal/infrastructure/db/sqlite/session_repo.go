package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.Session, error) {
	var s entities.Session
	var idStr, userIDStr, expiresAt, createdAt string
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?`, id.String()).
		Scan(&idStr, &userIDStr, &expiresAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying session: %w", err)
	}
	s.ID = uuid.MustParse(idStr)
	s.UserID = uuid.MustParse(userIDStr)
	s.ExpiresAt, _ = time.Parse(time.RFC3339, expiresAt)
	s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &s, nil
}

func (r *SessionRepo) Create(ctx context.Context, s *entities.Session) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`,
		s.ID.String(), s.UserID.String(), s.ExpiresAt.Format(time.RFC3339), s.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting session: %w", err)
	}
	return nil
}

func (r *SessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}

func (r *SessionRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = ?`, userID.String())
	if err != nil {
		return fmt.Errorf("deleting user sessions: %w", err)
	}
	return nil
}

func (r *SessionRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < ?`, time.Now().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("deleting expired sessions: %w", err)
	}
	return nil
}
