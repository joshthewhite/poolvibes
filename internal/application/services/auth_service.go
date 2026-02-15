package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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
	if len(cmd.Password) < 12 {
		return nil, nil, fmt.Errorf("password must be at least 12 characters")
	}
	if isCommonPassword(cmd.Password) {
		return nil, nil, fmt.Errorf("password is too common, please choose a stronger one")
	}

	existing, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("checking email: %w", err)
	}
	if existing != nil {
		return nil, nil, fmt.Errorf("unable to create account")
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
			slog.Warn("Failed to seed demo data", "userID", user.ID, "error", err)
		}
	}

	session := entities.NewSession(user.ID, sessionDuration)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("creating session: %w", err)
	}

	return user, session, nil
}

func (s *AuthService) SignIn(ctx context.Context, cmd command.SignIn) (*entities.User, *entities.Session, error) {
	// Use a dummy hash to compare against when user is not found,
	// so the timing is consistent regardless of whether the user exists.
	const dummyHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	user, err := s.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(cmd.Password))
		return nil, nil, fmt.Errorf("invalid email or password")
	}
	if user.IsDisabled {
		return nil, nil, fmt.Errorf("invalid email or password")
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

var commonPasswords = map[string]struct{}{
	"password":        {},
	"123456789012":    {},
	"qwertyuiop12":    {},
	"password1234":    {},
	"iloveyou1234":    {},
	"letmein12345":    {},
	"welcome12345":    {},
	"monkey123456":    {},
	"dragon123456":    {},
	"master123456":    {},
	"qwerty123456":    {},
	"login1234567":    {},
	"abc123456789":    {},
	"admin1234567":    {},
	"passw0rd1234":    {},
	"password12345":   {},
	"123456789abc":    {},
	"changeme1234":    {},
	"trustno11234":    {},
	"baseball1234":    {},
	"shadow123456":    {},
	"michael12345":    {},
	"football1234":    {},
	"superman1234":    {},
	"password1":       {},
	"password123":     {},
	"password1234567": {},
	"qwerty12345678":  {},
	"aaaaaaaaaaaa":    {},
	"123456789000":    {},
	"111111111111":    {},
	"000000000000":    {},
	"123123123123":    {},
	"abcdefghijkl":    {},
	"qwertyuiopas":    {},
}

func isCommonPassword(password string) bool {
	_, ok := commonPasswords[strings.ToLower(password)]
	return ok
}
