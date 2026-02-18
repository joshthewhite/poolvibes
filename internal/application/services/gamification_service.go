package services

import (
	"time"

	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

// ComputeHealthScore returns a 0-100 pool health score.
// Components: testing consistency (30%), water quality (30%), task completion (25%), chemical stock (15%).
func ComputeHealthScore(logs []entities.ChemistryLog, tasks []entities.Task, chemicals []entities.Chemical, now time.Time) int {
	if len(logs) == 0 && len(tasks) == 0 && len(chemicals) == 0 {
		return 0
	}

	// Testing Consistency (30%): tests in last 14 days / 4 (expected)
	fourteenDaysAgo := now.AddDate(0, 0, -14)
	recentTests := 0
	for _, l := range logs {
		if l.TestedAt.After(fourteenDaysAgo) {
			recentTests++
		}
	}
	testPct := float64(recentTests) / 4.0
	if testPct > 1 {
		testPct = 1
	}

	// Water Quality (30%): % of 6 readings in range on most recent test
	qualityPct := 0.0
	if len(logs) > 0 {
		latest := logs[0] // logs are newest-first
		inRange := 0
		if latest.PHInRange() {
			inRange++
		}
		if latest.FreeChlorineInRange() {
			inRange++
		}
		if latest.CombinedChlorineInRange() {
			inRange++
		}
		if latest.TotalAlkalinityInRange() {
			inRange++
		}
		if latest.CYAInRange() {
			inRange++
		}
		if latest.CalciumHardnessInRange() {
			inRange++
		}
		qualityPct = float64(inRange) / 6.0
	}

	// Task Completion (25%): % of tasks completed on time in last 30 days
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	totalDue := 0
	completedOnTime := 0
	for _, t := range tasks {
		if t.DueDate.Before(thirtyDaysAgo) {
			continue
		}
		if t.DueDate.After(now) {
			continue // not due yet
		}
		totalDue++
		if t.Status == entities.TaskStatusCompleted {
			completedOnTime++
		}
	}
	taskPct := 1.0 // default to 100% if no tasks were due
	if totalDue > 0 {
		taskPct = float64(completedOnTime) / float64(totalDue)
	}

	// Chemical Stock (15%): % of chemicals above threshold
	stockPct := 1.0 // default to 100% if no chemicals tracked
	if len(chemicals) > 0 {
		aboveThreshold := 0
		for _, c := range chemicals {
			if !c.IsLowStock() {
				aboveThreshold++
			}
		}
		stockPct = float64(aboveThreshold) / float64(len(chemicals))
	}

	score := int(testPct*30 + qualityPct*30 + taskPct*25 + stockPct*15)
	if score > 100 {
		score = 100
	}
	return score
}

// ComputeTestingStreak returns consecutive weeks with at least one water test.
func ComputeTestingStreak(logs []entities.ChemistryLog, now time.Time) int {
	if len(logs) == 0 {
		return 0
	}

	streak := 0
	for week := 0; ; week++ {
		weekEnd := now.AddDate(0, 0, -week*7)
		weekStart := now.AddDate(0, 0, -(week+1)*7)

		hasTest := false
		for _, l := range logs {
			if l.TestedAt.After(weekStart) && !l.TestedAt.After(weekEnd) {
				hasTest = true
				break
			}
		}
		if !hasTest {
			break
		}
		streak++
	}
	return streak
}

// ComputeTaskStreak returns consecutive weeks with zero overdue tasks.
func ComputeTaskStreak(tasks []entities.Task, now time.Time) int {
	streak := 0
	for week := 0; ; week++ {
		weekEnd := now.AddDate(0, 0, -week*7)
		weekStart := now.AddDate(0, 0, -(week+1)*7)

		hadOverdue := false
		for _, t := range tasks {
			if t.DueDate.Before(weekEnd) && t.DueDate.After(weekStart) {
				completedInTime := t.Status == entities.TaskStatusCompleted &&
					t.CompletedAt != nil && !t.CompletedAt.After(t.DueDate)
				if !completedInTime {
					hadOverdue = true
					break
				}
			}
		}
		if hadOverdue {
			break
		}
		streak++
		if streak >= 52 {
			break
		}
	}
	return streak
}

// CheckMilestones returns milestone keys newly earned (not in alreadyEarned).
func CheckMilestones(
	logs []entities.ChemistryLog,
	tasks []entities.Task,
	chemicals []entities.Chemical,
	healthScore int,
	alreadyEarned map[entities.MilestoneKey]bool,
) []entities.MilestoneKey {
	if alreadyEarned == nil {
		alreadyEarned = make(map[entities.MilestoneKey]bool)
	}

	var newly []entities.MilestoneKey
	check := func(key entities.MilestoneKey, met bool) {
		if met && !alreadyEarned[key] {
			newly = append(newly, key)
		}
	}

	// First Dip: at least one chemistry log
	check(entities.MilestoneFirstDip, len(logs) > 0)

	// Balanced: any test with all 6 readings in range
	balanced := false
	for _, l := range logs {
		if l.PHInRange() && l.FreeChlorineInRange() && l.CombinedChlorineInRange() &&
			l.TotalAlkalinityInRange() && l.CYAInRange() && l.CalciumHardnessInRange() {
			balanced = true
			break
		}
	}
	check(entities.MilestoneBalanced, balanced)

	// Consistent: 4-week testing streak
	now := time.Now()
	check(entities.MilestoneConsistent, ComputeTestingStreak(logs, now) >= 4)

	// Devoted: 12-week testing streak
	check(entities.MilestoneDevoted, ComputeTestingStreak(logs, now) >= 12)

	// On It: 10 tasks completed on time
	onTimeCount := 0
	for _, t := range tasks {
		if t.Status == entities.TaskStatusCompleted && t.CompletedAt != nil &&
			!t.CompletedAt.After(t.DueDate.AddDate(0, 0, 1)) {
			onTimeCount++
		}
	}
	check(entities.MilestoneOnIt, onTimeCount >= 10)

	// Stocked Up: all chemicals above threshold at once
	allStocked := len(chemicals) > 0
	for _, c := range chemicals {
		if c.IsLowStock() {
			allStocked = false
			break
		}
	}
	check(entities.MilestoneStockedUp, allStocked)

	// Clean Record: 30 days with zero overdue tasks (4+ week task streak)
	check(entities.MilestoneCleanRecord, ComputeTaskStreak(tasks, now) >= 4)

	// Pool Pro: health score >= 90
	check(entities.MilestonePoolPro, healthScore >= 90)

	return newly
}
