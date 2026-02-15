package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
	"golang.org/x/crypto/bcrypt"
)

const sessionDuration = 7 * 24 * time.Hour

type AuthService struct {
	userRepo     repositories.UserRepository
	sessionRepo  repositories.SessionRepository
	demoMode     bool
	maxDemoUsers int
	demoSeedSvc  *DemoSeedService
}

func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository, demoMode bool, maxDemoUsers int, demoSeedSvc *DemoSeedService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		demoMode:     demoMode,
		maxDemoUsers: maxDemoUsers,
		demoSeedSvc:  demoSeedSvc,
	}
}

func (s *AuthService) SignUp(ctx context.Context, cmd command.SignUp) (*entities.User, *entities.Session, error) {
	if cmd.Password == "" {
		return nil, nil, fmt.Errorf("password is required")
	}
	if len(cmd.Password) < 8 {
		return nil, nil, fmt.Errorf("password must be at least 8 characters")
	}

	existing, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("checking email: %w", err)
	}
	if existing != nil {
		return nil, nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("hashing password: %w", err)
	}

	user := entities.NewUser(cmd.Email, string(hash))

	// First user becomes admin automatically
	allUsers, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("checking users: %w", err)
	}
	if len(allUsers) == 0 {
		user.IsAdmin = true
	}

	// Non-admin users in demo mode get flagged as demo with 24h expiry
	if s.demoMode && !user.IsAdmin {
		if s.maxDemoUsers > 0 {
			count, err := s.userRepo.CountDemo(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("checking demo capacity: %w", err)
			}
			if count >= s.maxDemoUsers {
				return nil, nil, fmt.Errorf("demo slots are full, please try again later")
			}
		}
		user.IsDemo = true
		expires := time.Now().Add(24 * time.Hour)
		user.DemoExpiresAt = &expires
	}

	if err := user.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validation: %w", err)
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("creating user: %w", err)
	}

	// Seed demo data for demo users
	if user.IsDemo && s.demoSeedSvc != nil {
		if err := s.demoSeedSvc.Seed(ctx, user.ID); err != nil {
			log.Printf("Warning: failed to seed demo data for user %s: %v", user.ID, err)
		}
	}

	session := entities.NewSession(user.ID, sessionDuration)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("creating session: %w", err)
	}

	return user, session, nil
}

func (s *AuthService) SignIn(ctx context.Context, cmd command.SignIn) (*entities.User, *entities.Session, error) {
	user, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, nil, fmt.Errorf("invalid email or password")
	}
	if user.IsDisabled {
		return nil, nil, fmt.Errorf("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cmd.Password)); err != nil {
		return nil, nil, fmt.Errorf("invalid email or password")
	}

	session := entities.NewSession(user.ID, sessionDuration)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("creating session: %w", err)
	}

	return user, session, nil
}

func (s *AuthService) SignOut(ctx context.Context, sessionID string) error {
	uid, err := uuid.Parse(sessionID)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}
	return s.sessionRepo.Delete(ctx, uid)
}

func (s *AuthService) GetUserBySession(ctx context.Context, sessionID string) (*entities.User, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session ID: %w", err)
	}

	session, err := s.sessionRepo.FindByID(ctx, sid)
	if err != nil {
		return nil, fmt.Errorf("finding session: %w", err)
	}
	if session == nil || session.IsExpired() {
		if session != nil {
			_ = s.sessionRepo.Delete(ctx, sid)
		}
		return nil, nil
	}

	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user != nil && user.IsDisabled {
		return nil, nil
	}
	return user, nil
}
