package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type PushSubscriptionRepo struct {
	db *sql.DB
}

func NewPushSubscriptionRepo(db *sql.DB) *PushSubscriptionRepo {
	return &PushSubscriptionRepo{db: db}
}

func (r *PushSubscriptionRepo) Save(ctx context.Context, sub *entities.PushSubscription) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO push_subscriptions (id, user_id, endpoint, p256dh, auth, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, endpoint) DO UPDATE SET
			p256dh = excluded.p256dh,
			auth = excluded.auth`,
		sub.ID.String(), sub.UserID.String(), sub.Endpoint, sub.P256dh, sub.Auth,
		sub.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("saving push subscription: %w", err)
	}
	return nil
}

func (r *PushSubscriptionRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.PushSubscription, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, endpoint, p256dh, auth, created_at
		FROM push_subscriptions
		WHERE user_id = ?`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("querying push subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []entities.PushSubscription
	for rows.Next() {
		var s entities.PushSubscription
		var idStr, userIDStr, createdAt string
		if err := rows.Scan(&idStr, &userIDStr, &s.Endpoint, &s.P256dh, &s.Auth, &createdAt); err != nil {
			return nil, fmt.Errorf("scanning push subscription: %w", err)
		}
		s.ID = uuid.MustParse(idStr)
		s.UserID = uuid.MustParse(userIDStr)
		s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

func (r *PushSubscriptionRepo) DeleteByEndpoint(ctx context.Context, userID uuid.UUID, endpoint string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM push_subscriptions WHERE user_id = ? AND endpoint = ?`,
		userID.String(), endpoint)
	if err != nil {
		return fmt.Errorf("deleting push subscription: %w", err)
	}
	return nil
}

func (r *PushSubscriptionRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM push_subscriptions WHERE user_id = ?`, userID.String())
	if err != nil {
		return fmt.Errorf("deleting push subscriptions: %w", err)
	}
	return nil
}
