package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewMilestone(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	m := NewMilestone(userID, MilestoneFirstDip)

	if m.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if m.UserID != userID {
		t.Error("expected matching userID")
	}
	if m.Milestone != MilestoneFirstDip {
		t.Errorf("expected MilestoneFirstDip, got %s", m.Milestone)
	}
	if m.EarnedAt.IsZero() {
		t.Error("expected non-zero EarnedAt")
	}
}

func TestMilestoneKey_Valid(t *testing.T) {
	tests := []struct {
		key  MilestoneKey
		want bool
	}{
		{MilestoneFirstDip, true},
		{MilestoneBalanced, true},
		{MilestoneConsistent, true},
		{MilestoneDevoted, true},
		{MilestoneOnIt, true},
		{MilestoneStockedUp, true},
		{MilestoneCleanRecord, true},
		{MilestonePoolPro, true},
		{MilestoneKey("invalid"), false},
	}
	for _, tt := range tests {
		t.Run(string(tt.key), func(t *testing.T) {
			if got := tt.key.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllMilestones(t *testing.T) {
	all := AllMilestones()
	if len(all) != 8 {
		t.Errorf("expected 8 milestones, got %d", len(all))
	}
}
