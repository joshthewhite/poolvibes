package postgres

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
	rows, err := r.db.QueryContext(ctx, `SELECT id, email, password_hash, is_admin, is_disabled, created_at, updated_at FROM users ORDER BY created_at DESC`)
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
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, is_admin, is_disabled, created_at, updated_at FROM users WHERE id = $1`, id)
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
	row := r.db.QueryRowContext(ctx, `SELECT id, email, password_hash, is_admin, is_disabled, created_at, updated_at FROM users WHERE email = $1`, email)
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
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, email, password_hash, is_admin, is_disabled, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID, u.Email, u.PasswordHash, u.IsAdmin, u.IsDisabled, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, u *entities.User) error {
	u.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE users SET email = $1, password_hash = $2, is_admin = $3, is_disabled = $4, updated_at = $5 WHERE id = $6`,
		u.Email, u.PasswordHash, u.IsAdmin, u.IsDisabled, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUserFromRow(s scanner) (*entities.User, error) {
	var u entities.User
	if err := s.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.IsDisabled, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func scanUser(rows *sql.Rows) (*entities.User, error) {
	return scanUserFromRow(rows)
}

func scanUserRow(row *sql.Row) (*entities.User, error) {
	return scanUserFromRow(row)
}
