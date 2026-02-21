package postgres

import (
	"context"
	"database/sql"
	"fmt"

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
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT(user_id, endpoint) DO UPDATE SET
			p256dh = EXCLUDED.p256dh,
			auth = EXCLUDED.auth`,
		sub.ID, sub.UserID, sub.Endpoint, sub.P256dh, sub.Auth, sub.CreatedAt)
	if err != nil {
		return fmt.Errorf("saving push subscription: %w", err)
	}
	return nil
}

func (r *PushSubscriptionRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.PushSubscription, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, endpoint, p256dh, auth, created_at
		FROM push_subscriptions
		WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying push subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []entities.PushSubscription
	for rows.Next() {
		var s entities.PushSubscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.Endpoint, &s.P256dh, &s.Auth, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning push subscription: %w", err)
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

func (r *PushSubscriptionRepo) DeleteByEndpoint(ctx context.Context, userID uuid.UUID, endpoint string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM push_subscriptions WHERE user_id = $1 AND endpoint = $2`,
		userID, endpoint)
	if err != nil {
		return fmt.Errorf("deleting push subscription: %w", err)
	}
	return nil
}

func (r *PushSubscriptionRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM push_subscriptions WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("deleting push subscriptions: %w", err)
	}
	return nil
}
