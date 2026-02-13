package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type contextKey string

const userContextKey contextKey = "user"

func WithUser(ctx context.Context, user *entities.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (*entities.User, error) {
	user, ok := ctx.Value(userContextKey).(*entities.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	user, err := UserFromContext(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}
