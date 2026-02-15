package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"golang.org/x/crypto/bcrypt"
)

// --- mock repos ---

type mockUserRepo struct {
	users []*entities.User
}

func (m *mockUserRepo) FindAll(_ context.Context) ([]entities.User, error) {
	out := make([]entities.User, len(m.users))
	for i, u := range m.users {
		out[i] = *u
	}
	return out, nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id uuid.UUID) (*entities.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email string) (*entities.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) Create(_ context.Context, user *entities.User) error {
	m.users = append(m.users, user)
	return nil
}

func (m *mockUserRepo) Update(_ context.Context, user *entities.User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			m.users[i] = user
			return nil
		}
	}
	return nil
}

func (m *mockUserRepo) Delete(_ context.Context, id uuid.UUID) error {
	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockUserRepo) FindExpiredDemo(_ context.Context, now time.Time) ([]entities.User, error) {
	var result []entities.User
	for _, u := range m.users {
		if u.IsDemo && u.DemoExpiresAt != nil && now.After(*u.DemoExpiresAt) {
			result = append(result, *u)
		}
	}
	return result, nil
}

type mockSessionRepo struct {
	sessions []*entities.Session
}

func (m *mockSessionRepo) FindByID(_ context.Context, id uuid.UUID) (*entities.Session, error) {
	for _, s := range m.sessions {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, nil
}

func (m *mockSessionRepo) Create(_ context.Context, session *entities.Session) error {
	m.sessions = append(m.sessions, session)
	return nil
}

func (m *mockSessionRepo) Delete(_ context.Context, id uuid.UUID) error {
	for i, s := range m.sessions {
		if s.ID == id {
			m.sessions = append(m.sessions[:i], m.sessions[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockSessionRepo) DeleteByUserID(_ context.Context, userID uuid.UUID) error {
	filtered := m.sessions[:0]
	for _, s := range m.sessions {
		if s.UserID != userID {
			filtered = append(filtered, s)
		}
	}
	m.sessions = filtered
	return nil
}

func (m *mockSessionRepo) DeleteExpired(_ context.Context) error {
	filtered := m.sessions[:0]
	for _, s := range m.sessions {
		if !s.IsExpired() {
			filtered = append(filtered, s)
		}
	}
	m.sessions = filtered
	return nil
}

// --- SignUp tests ---

func TestAuthService_SignUp(t *testing.T) {
	tests := []struct {
		name        string
		cmd         command.SignUp
		existingCnt int // number of pre-existing users
		wantErr     string
		wantAdmin   bool
	}{
		{
			name:    "empty password",
			cmd:     command.SignUp{Email: "a@b.com", Password: ""},
			wantErr: "password is required",
		},
		{
			name:    "short password",
			cmd:     command.SignUp{Email: "a@b.com", Password: "short"},
			wantErr: "password must be at least 8 characters",
		},
		{
			name:      "first user becomes admin",
			cmd:       command.SignUp{Email: "admin@pool.com", Password: "password123"},
			wantAdmin: true,
		},
		{
			name:        "second user is not admin",
			cmd:         command.SignUp{Email: "user@pool.com", Password: "password123"},
			existingCnt: 1,
			wantAdmin:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mockUserRepo{}
			sessionRepo := &mockSessionRepo{}
			for i := 0; i < tt.existingCnt; i++ {
				userRepo.users = append(userRepo.users, &entities.User{
					ID:           uuid.New(),
					Email:        "existing@pool.com",
					PasswordHash: "hash",
				})
			}
			svc := NewAuthService(userRepo, sessionRepo, false, nil)

			user, session, err := svc.SignUp(context.Background(), tt.cmd)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantErr {
					t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user == nil {
				t.Fatal("expected user, got nil")
			}
			if session == nil {
				t.Fatal("expected session, got nil")
			}
			if user.IsAdmin != tt.wantAdmin {
				t.Errorf("IsAdmin = %v, want %v", user.IsAdmin, tt.wantAdmin)
			}
			if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.cmd.Password)); err != nil {
				t.Error("password hash does not match")
			}
			if session.UserID != user.ID {
				t.Errorf("session UserID = %v, want %v", session.UserID, user.ID)
			}
		})
	}
}

func TestAuthService_SignUp_DuplicateEmail(t *testing.T) {
	userRepo := &mockUserRepo{
		users: []*entities.User{{ID: uuid.New(), Email: "dup@pool.com", PasswordHash: "hash"}},
	}
	svc := NewAuthService(userRepo, &mockSessionRepo{}, false, nil)

	_, _, err := svc.SignUp(context.Background(), command.SignUp{
		Email:    "dup@pool.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	if err.Error() != "email already registered" {
		t.Errorf("error = %q, want %q", err.Error(), "email already registered")
	}
}

// --- SignIn tests ---

func TestAuthService_SignIn(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	activeUser := &entities.User{
		ID:           uuid.New(),
		Email:        "user@pool.com",
		PasswordHash: string(hash),
	}
	disabledUser := &entities.User{
		ID:           uuid.New(),
		Email:        "disabled@pool.com",
		PasswordHash: string(hash),
		IsDisabled:   true,
	}

	tests := []struct {
		name    string
		cmd     command.SignIn
		wantErr string
	}{
		{
			name: "valid credentials",
			cmd:  command.SignIn{Email: "user@pool.com", Password: "password123"},
		},
		{
			name:    "unknown email",
			cmd:     command.SignIn{Email: "unknown@pool.com", Password: "password123"},
			wantErr: "invalid email or password",
		},
		{
			name:    "wrong password",
			cmd:     command.SignIn{Email: "user@pool.com", Password: "wrongpass"},
			wantErr: "invalid email or password",
		},
		{
			name:    "disabled account",
			cmd:     command.SignIn{Email: "disabled@pool.com", Password: "password123"},
			wantErr: "account is disabled",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mockUserRepo{users: []*entities.User{activeUser, disabledUser}}
			sessionRepo := &mockSessionRepo{}
			svc := NewAuthService(userRepo, sessionRepo, false, nil)

			user, session, err := svc.SignIn(context.Background(), tt.cmd)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantErr {
					t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user == nil || session == nil {
				t.Fatal("expected user and session")
			}
			if session.UserID != user.ID {
				t.Errorf("session UserID = %v, want %v", session.UserID, user.ID)
			}
		})
	}
}

// --- SignOut tests ---

func TestAuthService_SignOut(t *testing.T) {
	sessionID := uuid.New()
	sessionRepo := &mockSessionRepo{
		sessions: []*entities.Session{{ID: sessionID, UserID: uuid.New(), ExpiresAt: time.Now().Add(time.Hour)}},
	}
	svc := NewAuthService(&mockUserRepo{}, sessionRepo, false, nil)

	if err := svc.SignOut(context.Background(), sessionID.String()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sessionRepo.sessions) != 0 {
		t.Error("session should have been deleted")
	}
}

func TestAuthService_SignOut_InvalidID(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{}, &mockSessionRepo{}, false, nil)
	if err := svc.SignOut(context.Background(), "not-a-uuid"); err == nil {
		t.Fatal("expected error for invalid session ID")
	}
}

// --- GetUserBySession tests ---

func TestAuthService_GetUserBySession(t *testing.T) {
	userID := uuid.New()
	user := &entities.User{ID: userID, Email: "u@pool.com", PasswordHash: "hash"}
	disabledUser := &entities.User{ID: uuid.New(), Email: "d@pool.com", PasswordHash: "hash", IsDisabled: true}

	validSession := &entities.Session{ID: uuid.New(), UserID: userID, ExpiresAt: time.Now().Add(time.Hour)}
	expiredSession := &entities.Session{ID: uuid.New(), UserID: userID, ExpiresAt: time.Now().Add(-time.Hour)}
	disabledSession := &entities.Session{ID: uuid.New(), UserID: disabledUser.ID, ExpiresAt: time.Now().Add(time.Hour)}

	tests := []struct {
		name     string
		sid      string
		wantUser bool
		wantErr  bool
	}{
		{"valid session", validSession.ID.String(), true, false},
		{"expired session", expiredSession.ID.String(), false, false},
		{"unknown session", uuid.New().String(), false, false},
		{"disabled user session", disabledSession.ID.String(), false, false},
		{"invalid uuid", "bad", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mockUserRepo{users: []*entities.User{user, disabledUser}}
			sessionRepo := &mockSessionRepo{sessions: []*entities.Session{validSession, expiredSession, disabledSession}}
			svc := NewAuthService(userRepo, sessionRepo, false, nil)

			u, err := svc.GetUserBySession(context.Background(), tt.sid)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantUser && u == nil {
				t.Fatal("expected user, got nil")
			}
			if !tt.wantUser && u != nil {
				t.Errorf("expected nil user, got %v", u.Email)
			}
		})
	}
}
