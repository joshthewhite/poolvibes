package entities

import "testing"

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr string
	}{
		{
			name:    "empty email",
			user:    User{Email: "", PasswordHash: "hash"},
			wantErr: "email is required",
		},
		{
			name:    "email without @",
			user:    User{Email: "invalid", PasswordHash: "hash"},
			wantErr: "email is invalid",
		},
		{
			name:    "empty password hash",
			user:    User{Email: "user@example.com", PasswordHash: ""},
			wantErr: "password hash is required",
		},
		{
			name: "valid",
			user: User{Email: "user@example.com", PasswordHash: "hash"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
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
		})
	}
}
