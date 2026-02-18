package entities

import (
	"time"

	"github.com/google/uuid"
)

type MilestoneKey string

const (
	MilestoneFirstDip    MilestoneKey = "first_dip"
	MilestoneBalanced    MilestoneKey = "balanced"
	MilestoneConsistent  MilestoneKey = "consistent"
	MilestoneDevoted     MilestoneKey = "devoted"
	MilestoneOnIt        MilestoneKey = "on_it"
	MilestoneStockedUp   MilestoneKey = "stocked_up"
	MilestoneCleanRecord MilestoneKey = "clean_record"
	MilestonePoolPro     MilestoneKey = "pool_pro"
)

func AllMilestones() []MilestoneKey {
	return []MilestoneKey{
		MilestoneFirstDip, MilestoneBalanced, MilestoneConsistent, MilestoneDevoted,
		MilestoneOnIt, MilestoneStockedUp, MilestoneCleanRecord, MilestonePoolPro,
	}
}

func (k MilestoneKey) Valid() bool {
	for _, m := range AllMilestones() {
		if k == m {
			return true
		}
	}
	return false
}

type Milestone struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Milestone MilestoneKey
	EarnedAt  time.Time
}

func NewMilestone(userID uuid.UUID, key MilestoneKey) *Milestone {
	return &Milestone{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    userID,
		Milestone: key,
		EarnedAt:  time.Now(),
	}
}
