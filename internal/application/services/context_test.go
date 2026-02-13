package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

func TestWithUser_UserFromContext(t *testing.T) {
	user := &entities.User{ID: uuid.New(), Email: "test@pool.com"}
	ctx := WithUser(context.Background(), user)

	got, err := UserFromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("got user ID %v, want %v", got.ID, user.ID)
	}
}

func TestUserFromContext_Missing(t *testing.T) {
	_, err := UserFromContext(context.Background())
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}

func TestUserIDFromContext(t *testing.T) {
	user := &entities.User{ID: uuid.New(), Email: "test@pool.com"}
	ctx := WithUser(context.Background(), user)

	id, err := UserIDFromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != user.ID {
		t.Errorf("got ID %v, want %v", id, user.ID)
	}
}

func TestUserIDFromContext_Missing(t *testing.T) {
	_, err := UserIDFromContext(context.Background())
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}
