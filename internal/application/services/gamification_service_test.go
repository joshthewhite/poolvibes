package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

func stockQty(amount float64) valueobjects.Quantity {
	q, _ := valueobjects.NewQuantity(amount, valueobjects.UnitPounds)
	return q
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestComputeHealthScore_AllPerfect(t *testing.T) {
	now := time.Now()
	userID := uuid.Must(uuid.NewV7())

	logs := make([]entities.ChemistryLog, 4)
	for i := range logs {
		logs[i] = entities.ChemistryLog{
			ID: uuid.Must(uuid.NewV7()), UserID: userID,
			PH: 7.4, FreeChlorine: 2.0, CombinedChlorine: 0.2,
			TotalAlkalinity: 100, CYA: 40, CalciumHardness: 300,
			TestedAt: now.AddDate(0, 0, -i*3),
		}
	}

	tasks := []entities.Task{
		{ID: uuid.Must(uuid.NewV7()), UserID: userID, Status: entities.TaskStatusCompleted,
			DueDate: now.AddDate(0, 0, -5)},
	}

	chemicals := []entities.Chemical{
		{ID: uuid.Must(uuid.NewV7()), UserID: userID, Stock: stockQty(10), AlertThreshold: 5},
	}

	score := ComputeHealthScore(logs, tasks, chemicals, now)
	if score < 90 {
		t.Errorf("expected score >= 90 for perfect data, got %d", score)
	}
}

func TestComputeHealthScore_NoData(t *testing.T) {
	score := ComputeHealthScore(nil, nil, nil, time.Now())
	if score != 0 {
		t.Errorf("expected 0 for no data, got %d", score)
	}
}

func TestComputeTestingStreak(t *testing.T) {
	now := time.Now()
	userID := uuid.Must(uuid.NewV7())

	logs := []entities.ChemistryLog{
		{UserID: userID, TestedAt: now.AddDate(0, 0, -1)},
		{UserID: userID, TestedAt: now.AddDate(0, 0, -8)},
		{UserID: userID, TestedAt: now.AddDate(0, 0, -15)},
	}

	streak := ComputeTestingStreak(logs, now)
	if streak != 3 {
		t.Errorf("expected 3-week testing streak, got %d", streak)
	}
}

func TestComputeTestingStreak_NoLogs(t *testing.T) {
	streak := ComputeTestingStreak(nil, time.Now())
	if streak != 0 {
		t.Errorf("expected 0, got %d", streak)
	}
}

func TestComputeTaskStreak(t *testing.T) {
	now := time.Now()
	userID := uuid.Must(uuid.NewV7())

	tasks := []entities.Task{
		{UserID: userID, Status: entities.TaskStatusCompleted,
			DueDate: now.AddDate(0, 0, -2), CompletedAt: timePtr(now.AddDate(0, 0, -2))},
		{UserID: userID, Status: entities.TaskStatusCompleted,
			DueDate: now.AddDate(0, 0, -9), CompletedAt: timePtr(now.AddDate(0, 0, -9))},
		{UserID: userID, Status: entities.TaskStatusCompleted,
			DueDate: now.AddDate(0, 0, -16), CompletedAt: timePtr(now.AddDate(0, 0, -16))},
	}

	streak := ComputeTaskStreak(tasks, now)
	if streak < 3 {
		t.Errorf("expected >= 3-week task streak, got %d", streak)
	}
}

func TestCheckMilestones_FirstDip(t *testing.T) {
	logs := []entities.ChemistryLog{
		{ID: uuid.Must(uuid.NewV7()), TestedAt: time.Now()},
	}
	earned := CheckMilestones(logs, nil, nil, 0, nil)
	found := false
	for _, m := range earned {
		if m == entities.MilestoneFirstDip {
			found = true
		}
	}
	if !found {
		t.Error("expected MilestoneFirstDip to be earned")
	}
}

func TestCheckMilestones_Balanced(t *testing.T) {
	logs := []entities.ChemistryLog{
		{
			PH: 7.4, FreeChlorine: 2.0, CombinedChlorine: 0.2,
			TotalAlkalinity: 100, CYA: 40, CalciumHardness: 300,
			TestedAt: time.Now(),
		},
	}
	earned := CheckMilestones(logs, nil, nil, 0, nil)
	found := false
	for _, m := range earned {
		if m == entities.MilestoneBalanced {
			found = true
		}
	}
	if !found {
		t.Error("expected MilestoneBalanced to be earned")
	}
}

func TestCheckMilestones_SkipsAlreadyEarned(t *testing.T) {
	logs := []entities.ChemistryLog{
		{TestedAt: time.Now()},
	}
	alreadyEarned := map[entities.MilestoneKey]bool{entities.MilestoneFirstDip: true}
	earned := CheckMilestones(logs, nil, nil, 0, alreadyEarned)
	for _, m := range earned {
		if m == entities.MilestoneFirstDip {
			t.Error("should not re-earn MilestoneFirstDip")
		}
	}
}

func TestCheckMilestones_PoolPro(t *testing.T) {
	earned := CheckMilestones(nil, nil, nil, 92, nil)
	found := false
	for _, m := range earned {
		if m == entities.MilestonePoolPro {
			found = true
		}
	}
	if !found {
		t.Error("expected MilestonePoolPro for score 92")
	}
}
