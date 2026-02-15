package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

type UserService struct {
	repo        repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

func NewUserService(repo repositories.UserRepository, sessionRepo repositories.SessionRepository) *UserService {
	return &UserService{repo: repo, sessionRepo: sessionRepo}
}

func (s *UserService) List(ctx context.Context) ([]entities.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) Get(ctx context.Context, id string) (*entities.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.FindByID(ctx, uid)
}

func (s *UserService) Update(ctx context.Context, cmd command.UpdateUser) (*entities.User, error) {
	uid, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %w", err)
	}
	user, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	user.IsAdmin = cmd.IsAdmin
	user.IsDisabled = cmd.IsDisabled
	if !cmd.IsDemo && user.IsDemo {
		user.IsDemo = false
		user.DemoExpiresAt = nil
	}
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	if cmd.IsDisabled {
		_ = s.sessionRepo.DeleteByUserID(ctx, uid)
	}
	return user, nil
}

func (s *UserService) UpdatePreferences(ctx context.Context, cmd command.UpdateNotificationPreferences) (*entities.User, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	user.Phone = cmd.Phone
	user.NotifyEmail = cmd.NotifyEmail
	user.NotifySMS = cmd.NotifySMS
	user.PoolGallons = cmd.PoolGallons
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}
	return s.repo.Delete(ctx, uid)
}
