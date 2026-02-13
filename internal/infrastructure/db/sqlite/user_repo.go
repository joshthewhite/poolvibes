package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindAll(ctx context.Context) ([]entities.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, email, password_hash, is_admin, is_disabled,
			phone, notify_email, notify_sms,
			created_at, updated_at
		FROM users
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, is_admin, is_disabled,
			phone, notify_email, notify_sms,
			created_at, updated_at
		FROM users
		WHERE id = ?`, id.String())
	u, err := scanUserRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying user: %w", err)
	}
	return u, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, is_admin, is_disabled,
			phone, notify_email, notify_sms,
			created_at, updated_at
		FROM users
		WHERE email = ?`, email)
	u, err := scanUserRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by email: %w", err)
	}
	return u, nil
}

func (r *UserRepo) Create(ctx context.Context, u *entities.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, is_admin,
			is_disabled, phone, notify_email, notify_sms,
			created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.ID.String(), u.Email, u.PasswordHash, boolToInt(u.IsAdmin), boolToInt(u.IsDisabled),
		u.Phone, boolToInt(u.NotifyEmail), boolToInt(u.NotifySMS),
		u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, u *entities.User) error {
	u.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET email = ?, password_hash = ?,
			is_admin = ?, is_disabled = ?,
			phone = ?, notify_email = ?, notify_sms = ?,
			updated_at = ?
		WHERE id = ?`,
		u.Email, u.PasswordHash, boolToInt(u.IsAdmin), boolToInt(u.IsDisabled),
		u.Phone, boolToInt(u.NotifyEmail), boolToInt(u.NotifySMS),
		u.UpdatedAt.Format(time.RFC3339), u.ID.String())
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id.String())
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scanUserFromRow(s scanner) (*entities.User, error) {
	var u entities.User
	var idStr, createdAt, updatedAt string
	var isAdmin, isDisabled, notifyEmail, notifySMS int
	if err := s.Scan(&idStr, &u.Email, &u.PasswordHash, &isAdmin, &isDisabled,
		&u.Phone, &notifyEmail, &notifySMS,
		&createdAt, &updatedAt); err != nil {
		return nil, err
	}
	u.ID = uuid.MustParse(idStr)
	u.IsAdmin = isAdmin == 1
	u.IsDisabled = isDisabled == 1
	u.NotifyEmail = notifyEmail == 1
	u.NotifySMS = notifySMS == 1
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	u.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &u, nil
}

func scanUser(rows *sql.Rows) (*entities.User, error) {
	return scanUserFromRow(rows)
}

func scanUserRow(row *sql.Row) (*entities.User, error) {
	return scanUserFromRow(row)
}
