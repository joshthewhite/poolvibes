package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	IsAdmin      bool
	IsDisabled   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.Must(uuid.NewV7()),
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		IsAdmin:      false,
		IsDisabled:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("email is invalid")
	}
	if u.PasswordHash == "" {
		return fmt.Errorf("password hash is required")
	}
	return nil
}
