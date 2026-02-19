# Gamification Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Pool Health Score (0-100), streaks, and 8 milestone badges to the PoolVibes dashboard.

**Architecture:** A new `GamificationService` computes score, streaks, and checks milestones on each dashboard load. Score and streaks are computed on-the-fly from existing data. Milestones persist in a new `user_milestones` DB table. The dashboard handler calls the gamification service and passes data to new templ components. Font Awesome 6 Free CDN provides milestone icons.

**Tech Stack:** Go, templ, Bulma, Datastar, Font Awesome 6 Free CDN, SQLite + Postgres migrations

---

## Task 1: Database Migration — `user_milestones` Table

**Files:**
- Create: `migrations/sqlite/000007_user_milestones.up.sql`
- Create: `migrations/sqlite/000007_user_milestones.down.sql`
- Create: `migrations/postgres/000007_user_milestones.up.sql`
- Create: `migrations/postgres/000007_user_milestones.down.sql`

**Step 1: Create SQLite up migration**

```sql
CREATE TABLE IF NOT EXISTS user_milestones (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    milestone TEXT NOT NULL,
    earned_at TEXT NOT NULL,
    UNIQUE(user_id, milestone)
);
CREATE INDEX idx_user_milestones_user_id ON user_milestones(user_id);
```

**Step 2: Create SQLite down migration**

```sql
DROP TABLE IF EXISTS user_milestones;
```

**Step 3: Create Postgres up migration**

```sql
CREATE TABLE IF NOT EXISTS user_milestones (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    milestone TEXT NOT NULL,
    earned_at TIMESTAMPTZ NOT NULL,
    UNIQUE(user_id, milestone)
);
CREATE INDEX idx_user_milestones_user_id ON user_milestones(user_id);
```

**Step 4: Create Postgres down migration**

```sql
DROP TABLE IF EXISTS user_milestones;
```

**Step 5: Verify migration applies**

Run: `task build && ./bin/poolvibes serve --db /tmp/test-gamification.db`
Expected: Server starts without migration errors. Kill it after confirming.
Clean up: `rm /tmp/test-gamification.db`

**Step 6: Commit**

```bash
git add migrations/
git commit -m "feat: add user_milestones migration"
```

---

## Task 2: Milestone Entity

**Files:**
- Create: `internal/domain/entities/milestone.go`
- Create: `internal/domain/entities/milestone_test.go`

**Step 1: Write the test**

```go
package entities

import (
	"testing"
	"time"

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
```

**Step 2: Run test to verify it fails**

Run: `go test -v -run TestNewMilestone ./internal/domain/entities/`
Expected: FAIL — `NewMilestone` not defined

**Step 3: Write the entity**

```go
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
```

**Step 4: Run tests to verify they pass**

Run: `go test -v -run "TestNewMilestone|TestMilestoneKey|TestAllMilestones" ./internal/domain/entities/`
Expected: PASS (3 tests)

**Step 5: Commit**

```bash
git add internal/domain/entities/milestone.go internal/domain/entities/milestone_test.go
git commit -m "feat: add Milestone entity with MilestoneKey enum"
```

---

## Task 3: Milestone Repository Interface

**Files:**
- Create: `internal/domain/repositories/milestone.go`

**Step 1: Write the interface**

```go
package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type MilestoneRepository interface {
	FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Milestone, error)
	Create(ctx context.Context, milestone *entities.Milestone) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
```

Note: `DeleteByUserID` is needed for demo user cleanup. No `FindByID`, `Update`, or `Delete` by ID — milestones are write-once, read-many.

**Step 2: Verify it compiles**

Run: `go build ./internal/domain/repositories/`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/domain/repositories/milestone.go
git commit -m "feat: add MilestoneRepository interface"
```

---

## Task 4: SQLite Milestone Repository

**Files:**
- Create: `internal/infrastructure/db/sqlite/milestone_repo.go`
- Create: `internal/infrastructure/db/sqlite/milestone_repo_test.go`

**Step 1: Write the test**

Check test patterns in the project first. If there are no integration tests using real SQLite, use a simple compile-and-interface-satisfaction test:

```go
package sqlite

import (
	"testing"

	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

func TestMilestoneRepoImplementsInterface(t *testing.T) {
	var _ repositories.MilestoneRepository = (*MilestoneRepo)(nil)
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v -run TestMilestoneRepoImplementsInterface ./internal/infrastructure/db/sqlite/`
Expected: FAIL — `MilestoneRepo` not defined

**Step 3: Write the repo**

```go
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type MilestoneRepo struct {
	db *sql.DB
}

func NewMilestoneRepo(db *sql.DB) *MilestoneRepo {
	return &MilestoneRepo{db: db}
}

func (r *MilestoneRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Milestone, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, milestone, earned_at
		FROM user_milestones
		WHERE user_id = ?
		ORDER BY earned_at ASC`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("querying milestones: %w", err)
	}
	defer rows.Close()

	var milestones []entities.Milestone
	for rows.Next() {
		var m entities.Milestone
		var idStr, userIDStr, earnedAt string
		var milestone string
		if err := rows.Scan(&idStr, &userIDStr, &milestone, &earnedAt); err != nil {
			return nil, fmt.Errorf("scanning milestone: %w", err)
		}
		m.ID = uuid.MustParse(idStr)
		m.UserID = uuid.MustParse(userIDStr)
		m.Milestone = entities.MilestoneKey(milestone)
		m.EarnedAt, _ = time.Parse(time.RFC3339, earnedAt)
		milestones = append(milestones, m)
	}
	return milestones, rows.Err()
}

func (r *MilestoneRepo) Create(ctx context.Context, m *entities.Milestone) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_milestones (id, user_id, milestone, earned_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (user_id, milestone) DO NOTHING`,
		m.ID.String(), m.UserID.String(), string(m.Milestone), m.EarnedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("inserting milestone: %w", err)
	}
	return nil
}

func (r *MilestoneRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_milestones WHERE user_id = ?`, userID.String())
	if err != nil {
		return fmt.Errorf("deleting milestones: %w", err)
	}
	return nil
}
```

Key pattern notes for SQLite:
- `?` placeholders
- `.String()` on all UUIDs
- `time.RFC3339` for timestamp storage/parsing
- `ON CONFLICT ... DO NOTHING` so re-earning is idempotent

**Step 4: Run test to verify it passes**

Run: `go test -v -run TestMilestoneRepoImplementsInterface ./internal/infrastructure/db/sqlite/`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/infrastructure/db/sqlite/milestone_repo.go internal/infrastructure/db/sqlite/milestone_repo_test.go
git commit -m "feat: add SQLite milestone repository"
```

---

## Task 5: Postgres Milestone Repository

**Files:**
- Create: `internal/infrastructure/db/postgres/milestone_repo.go`
- Create: `internal/infrastructure/db/postgres/milestone_repo_test.go`

**Step 1: Write the test**

```go
package postgres

import (
	"testing"

	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
)

func TestMilestoneRepoImplementsInterface(t *testing.T) {
	var _ repositories.MilestoneRepository = (*MilestoneRepo)(nil)
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v -run TestMilestoneRepoImplementsInterface ./internal/infrastructure/db/postgres/`
Expected: FAIL — `MilestoneRepo` not defined

**Step 3: Write the repo**

```go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

type MilestoneRepo struct {
	db *sql.DB
}

func NewMilestoneRepo(db *sql.DB) *MilestoneRepo {
	return &MilestoneRepo{db: db}
}

func (r *MilestoneRepo) FindAll(ctx context.Context, userID uuid.UUID) ([]entities.Milestone, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, milestone, earned_at
		FROM user_milestones
		WHERE user_id = $1
		ORDER BY earned_at ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("querying milestones: %w", err)
	}
	defer rows.Close()

	var milestones []entities.Milestone
	for rows.Next() {
		var m entities.Milestone
		var milestone string
		if err := rows.Scan(&m.ID, &m.UserID, &milestone, &m.EarnedAt); err != nil {
			return nil, fmt.Errorf("scanning milestone: %w", err)
		}
		m.Milestone = entities.MilestoneKey(milestone)
		milestones = append(milestones, m)
	}
	return milestones, rows.Err()
}

func (r *MilestoneRepo) Create(ctx context.Context, m *entities.Milestone) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_milestones (id, user_id, milestone, earned_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, milestone) DO NOTHING`,
		m.ID, m.UserID, string(m.Milestone), m.EarnedAt)
	if err != nil {
		return fmt.Errorf("inserting milestone: %w", err)
	}
	return nil
}

func (r *MilestoneRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM user_milestones WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("deleting milestones: %w", err)
	}
	return nil
}
```

Key pattern notes for Postgres:
- `$1, $2` placeholders
- Native UUID and timestamp scanning (no `.String()` or `time.Parse` needed)

**Step 4: Run test to verify it passes**

Run: `go test -v -run TestMilestoneRepoImplementsInterface ./internal/infrastructure/db/postgres/`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/infrastructure/db/postgres/milestone_repo.go internal/infrastructure/db/postgres/milestone_repo_test.go
git commit -m "feat: add Postgres milestone repository"
```

---

## Task 6: Gamification Service — Score & Streaks

This service computes the Pool Health Score and streaks from existing data. Milestones come in Task 7.

**Files:**
- Create: `internal/application/services/gamification_service.go`
- Create: `internal/application/services/gamification_service_test.go`

**Step 1: Write the tests**

```go
package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
)

func TestComputeHealthScore_AllPerfect(t *testing.T) {
	now := time.Now()
	userID := uuid.Must(uuid.NewV7())

	// 4 tests in last 14 days → testing consistency = 100%
	logs := make([]entities.ChemistryLog, 4)
	for i := range logs {
		logs[i] = entities.ChemistryLog{
			ID:               uuid.Must(uuid.NewV7()),
			UserID:           userID,
			PH:               7.4,
			FreeChlorine:     2.0,
			CombinedChlorine: 0.2,
			TotalAlkalinity:  100,
			CYA:              40,
			CalciumHardness:  300,
			TestedAt:         now.AddDate(0, 0, -i*3),
		}
	}

	// All tasks completed on time (none overdue in last 30 days)
	tasks := []entities.Task{
		{ID: uuid.Must(uuid.NewV7()), UserID: userID, Status: entities.TaskStatusCompleted,
			DueDate: now.AddDate(0, 0, -5)},
	}

	// All chemicals above threshold
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

	// Tests in each of the last 3 weeks
	logs := []entities.ChemistryLog{
		{UserID: userID, TestedAt: now.AddDate(0, 0, -1)},
		{UserID: userID, TestedAt: now.AddDate(0, 0, -8)},
		{UserID: userID, TestedAt: now.AddDate(0, 0, -15)},
		// Gap — no test in week 4
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

	// All tasks on time for last 3 weeks, then overdue task in week 4
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

func timePtr(t time.Time) *time.Time {
	return &t
}
```

Note: `stockQty` is a test helper — define it in the test file:

```go
import "github.com/joshthewhite/poolvibes/internal/domain/valueobjects"

func stockQty(amount float64) valueobjects.Quantity {
	q, _ := valueobjects.NewQuantity(amount, valueobjects.UnitPounds)
	return q
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -v -run "TestCompute" ./internal/application/services/`
Expected: FAIL — `ComputeHealthScore` not defined

**Step 3: Write the gamification service**

```go
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
			// A task was overdue during this week if:
			// - its due date falls in or before this week AND
			// - it wasn't completed before end of this week (or not completed at all)
			if t.DueDate.Before(weekEnd) && t.DueDate.After(weekStart.AddDate(0, 0, -90)) {
				completedInTime := t.Status == entities.TaskStatusCompleted &&
					t.CompletedAt != nil && !t.CompletedAt.After(t.DueDate)
				if !completedInTime && t.DueDate.Before(weekEnd) && t.DueDate.After(weekStart) {
					hadOverdue = true
					break
				}
			}
		}
		if hadOverdue {
			break
		}
		streak++
		// Don't go back more than 52 weeks
		if streak >= 52 {
			break
		}
	}
	return streak
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -v -run "TestCompute" ./internal/application/services/`
Expected: PASS (5 tests)

**Step 5: Commit**

```bash
git add internal/application/services/gamification_service.go internal/application/services/gamification_service_test.go
git commit -m "feat: add gamification score and streak computation"
```

---

## Task 7: Gamification Service — Milestone Checking

**Files:**
- Modify: `internal/application/services/gamification_service.go`
- Modify: `internal/application/services/gamification_service_test.go`

**Step 1: Write the tests**

Add to `gamification_service_test.go`:

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test -v -run "TestCheckMilestones" ./internal/application/services/`
Expected: FAIL — `CheckMilestones` not defined

**Step 3: Add CheckMilestones to gamification_service.go**

```go
// CheckMilestones returns milestone keys newly earned (not in alreadyEarned).
// Parameters:
//   - logs: chemistry logs (newest first)
//   - tasks: all tasks
//   - chemicals: all chemicals
//   - healthScore: pre-computed health score
//   - alreadyEarned: set of already-earned milestone keys
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
```

**Step 4: Run tests to verify they pass**

Run: `go test -v -run "TestCheckMilestones" ./internal/application/services/`
Expected: PASS (4 tests)

**Step 5: Commit**

```bash
git add internal/application/services/gamification_service.go internal/application/services/gamification_service_test.go
git commit -m "feat: add milestone checking logic"
```

---

## Task 8: Dashboard Types — Gamification Data

**Files:**
- Modify: `internal/interface/web/templates/dashboard_types.go`

**Step 1: Add gamification types**

Add these types to `dashboard_types.go` and add gamification fields to `DashboardData`:

```go
type HealthScoreSummary struct {
	Score  int
	Status string // "good", "warning", "danger"
	Label  string // description text
}

type StreaksSummary struct {
	TestingStreak int
	TaskStreak    int
}

type MilestoneBadge struct {
	Key     string
	Name    string
	Icon    string // Font Awesome class
	Earned  bool
	IsNew   bool // just earned this load
}
```

Add fields to `DashboardData` struct:

```go
HealthScore HealthScoreSummary
Streaks     StreaksSummary
Milestones  []MilestoneBadge
```

**Step 2: Verify it compiles**

Run: `go build ./internal/interface/web/templates/`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/interface/web/templates/dashboard_types.go
git commit -m "feat: add gamification types to dashboard data"
```

---

## Task 9: Dashboard Handler — Wire Gamification

**Files:**
- Modify: `internal/interface/web/handlers/dashboard.go`

**Step 1: Update DashboardHandler to accept milestone repo**

Add a milestone service/repo dependency and update `buildDashboardData` to compute gamification data. The handler needs the milestone repo to read existing milestones and persist newly earned ones.

Update the handler struct and constructor:

```go
type DashboardHandler struct {
	chemSvc       *services.ChemistryService
	taskSvc       *services.TaskService
	chemicSvc     *services.ChemicalService
	milestoneRepo repositories.MilestoneRepository
}

func NewDashboardHandler(chemSvc *services.ChemistryService, taskSvc *services.TaskService, chemicSvc *services.ChemicalService, milestoneRepo repositories.MilestoneRepository) *DashboardHandler {
	return &DashboardHandler{chemSvc: chemSvc, taskSvc: taskSvc, chemicSvc: chemicSvc, milestoneRepo: milestoneRepo}
}
```

Add import: `"github.com/joshthewhite/poolvibes/internal/domain/repositories"`

Update `Page` to compute gamification and call `buildDashboardData` with the new params:

```go
func (h *DashboardHandler) Page(w http.ResponseWriter, r *http.Request) {
	logs, _ := h.chemSvc.List(r.Context())
	tasks, _ := h.taskSvc.List(r.Context())
	chemicals, _ := h.chemicSvc.List(r.Context())

	data := buildDashboardData(logs, tasks, chemicals)

	// Gamification
	now := time.Now()
	score := services.ComputeHealthScore(logs, tasks, chemicals, now)
	testStreak := services.ComputeTestingStreak(logs, now)
	taskStreak := services.ComputeTaskStreak(tasks, now)

	data.HealthScore = templates.HealthScoreSummary{
		Score:  score,
		Status: healthScoreStatus(score),
		Label:  healthScoreLabel(score),
	}
	data.Streaks = templates.StreaksSummary{
		TestingStreak: testStreak,
		TaskStreak:    taskStreak,
	}

	// Milestones
	existingMilestones, _ := h.milestoneRepo.FindAll(r.Context(), services.MustUserIDFromContext(r.Context()))
	earnedSet := make(map[entities.MilestoneKey]bool)
	for _, m := range existingMilestones {
		earnedSet[m.Milestone] = true
	}

	newlyEarned := services.CheckMilestones(logs, tasks, chemicals, score, earnedSet)

	// Persist newly earned milestones
	for _, key := range newlyEarned {
		userID, _ := services.UserIDFromContext(r.Context())
		m := entities.NewMilestone(userID, key)
		_ = h.milestoneRepo.Create(r.Context(), m)
		earnedSet[key] = true
	}
	newlyEarnedSet := make(map[entities.MilestoneKey]bool)
	for _, key := range newlyEarned {
		newlyEarnedSet[key] = true
	}

	data.Milestones = buildMilestoneBadges(earnedSet, newlyEarnedSet)

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.Dashboard(data))
}
```

Add helper functions:

```go
func healthScoreStatus(score int) string {
	if score >= 80 {
		return "good"
	}
	if score >= 50 {
		return "warning"
	}
	return "danger"
}

func healthScoreLabel(score int) string {
	if score >= 80 {
		return "Your pool is in great shape"
	}
	if score >= 50 {
		return "A few things need attention"
	}
	return "Your pool needs some love"
}

func buildMilestoneBadges(earned, newlyEarned map[entities.MilestoneKey]bool) []templates.MilestoneBadge {
	milestoneInfo := []struct {
		key  entities.MilestoneKey
		name string
		icon string
	}{
		{entities.MilestoneFirstDip, "First Dip", "fa-solid fa-droplet"},
		{entities.MilestoneBalanced, "Balanced", "fa-solid fa-scale-balanced"},
		{entities.MilestoneConsistent, "Consistent", "fa-solid fa-calendar-check"},
		{entities.MilestoneDevoted, "Devoted", "fa-solid fa-fire"},
		{entities.MilestoneOnIt, "On It", "fa-solid fa-clipboard-check"},
		{entities.MilestoneStockedUp, "Stocked Up", "fa-solid fa-boxes-stacked"},
		{entities.MilestoneCleanRecord, "Clean Record", "fa-solid fa-circle-check"},
		{entities.MilestonePoolPro, "Pool Pro", "fa-solid fa-trophy"},
	}

	badges := make([]templates.MilestoneBadge, len(milestoneInfo))
	for i, info := range milestoneInfo {
		badges[i] = templates.MilestoneBadge{
			Key:    string(info.key),
			Name:   info.name,
			Icon:   info.icon,
			Earned: earned[info.key],
			IsNew:  newlyEarned[info.key],
		}
	}
	return badges
}
```

Note: You may need to add a `MustUserIDFromContext` helper to `services/context.go`, or just use `UserIDFromContext` directly and ignore the error (since auth middleware guarantees the user exists). Check if it already exists; if not, use `UserIDFromContext` and handle the error.

**Step 2: Verify it compiles**

Run: `go build ./internal/interface/web/handlers/`
Expected: No errors (may need to fix imports — update `setupRoutes` call in server.go next task)

**Step 3: Commit**

```bash
git add internal/interface/web/handlers/dashboard.go
git commit -m "feat: wire gamification into dashboard handler"
```

---

## Task 10: Wiring — Serve + Server

**Files:**
- Modify: `cmd/serve.go` — add milestone repo variable and instantiation
- Modify: `internal/interface/web/server.go` — add milestone repo to Server, pass to DashboardHandler

**Step 1: Update `cmd/serve.go`**

Add `milestoneRepo repositories.MilestoneRepository` to the var block (line ~48).

In the `case "sqlite":` block, add:
```go
milestoneRepo = sqlite.NewMilestoneRepo(db)
```

In the `case "postgres":` block, add:
```go
milestoneRepo = postgres.NewMilestoneRepo(db)
```

Pass `milestoneRepo` to `web.NewServer`:
```go
server := web.NewServer(authSvc, userSvc, chemSvc, taskSvc, equipSvc, chemicSvc, milestoneRepo)
```

**Step 2: Update `internal/interface/web/server.go`**

Add `milestoneRepo repositories.MilestoneRepository` field to the `Server` struct.

Update `NewServer` signature to accept `milestoneRepo`:
```go
func NewServer(authSvc *services.AuthService, userSvc *services.UserService, chemSvc *services.ChemistryService, taskSvc *services.TaskService, equipSvc *services.EquipmentService, chemicSvc *services.ChemicalService, milestoneRepo repositories.MilestoneRepository) *Server
```

Add import: `"github.com/joshthewhite/poolvibes/internal/domain/repositories"`

Update `setupRoutes` to pass `milestoneRepo` to `NewDashboardHandler`:
```go
dashHandler := handlers.NewDashboardHandler(s.chemSvc, s.taskSvc, s.chemicSvc, s.milestoneRepo)
```

**Step 3: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 4: Commit**

```bash
git add cmd/serve.go internal/interface/web/server.go
git commit -m "feat: wire milestone repository through server"
```

---

## Task 11: Font Awesome CDN + Gamification CSS

**Files:**
- Modify: `internal/interface/web/templates/layout.templ`
- Modify: `internal/interface/web/server.go` (CSP header update)

**Step 1: Add Font Awesome CDN**

In `layout.templ`, add after the Bulma CSS link (line ~13):
```html
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css"/>
```

**Step 2: Update CSP header**

In `server.go`, update the `style-src` directive in `securityHeaders` to allow `cdnjs.cloudflare.com`:
```go
"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com https://cdnjs.cloudflare.com; "+
"font-src https://fonts.gstatic.com https://cdnjs.cloudflare.com; "+
```

**Step 3: Add gamification CSS**

In `layout.templ`, add these styles inside the existing `<style>` block, after the neumorphic styles:

```css
/* Pool Health Score */
.pv-health-score {
    text-align: center;
    padding: 1.5rem;
}
.pv-health-score .pv-score-number {
    font-size: 3.5rem;
    font-weight: 800;
    font-family: 'Inter Tight', sans-serif;
    line-height: 1;
}
.pv-health-score .pv-score-label {
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

/* Streaks */
.pv-streaks {
    text-align: center;
    font-size: 0.8rem;
    color: var(--pv-text-secondary);
    padding: 0.25rem 0 0.75rem;
}
.pv-streaks .pv-streak-item {
    display: inline;
}
.pv-streaks .pv-streak-item + .pv-streak-item::before {
    content: " · ";
}

/* Milestone badges */
.pv-milestones {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    justify-content: center;
    padding-bottom: 1rem;
}
.pv-milestone-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.35rem 0.75rem;
    border-radius: 999px;
    font-size: 0.75rem;
    font-weight: 500;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.pv-milestone-badge.is-earned {
    background: var(--pv-primary-light);
    color: var(--pv-primary);
}
.pv-milestone-badge.is-locked {
    background: #f0eef5;
    color: #b0adc0;
}
.pv-milestone-badge.is-new {
    animation: pv-milestone-glow 1.5s ease-in-out;
}
@keyframes pv-milestone-glow {
    0%, 100% { transform: scale(1); box-shadow: none; }
    50% { transform: scale(1.08); box-shadow: 0 0 12px rgba(13, 148, 136, 0.4); }
}
```

Add dark mode overrides inside the existing `@media (prefers-color-scheme: dark)` block:

```css
.pv-milestone-badge.is-earned {
    background: rgba(13, 148, 136, 0.15);
    color: var(--pv-primary);
}
.pv-milestone-badge.is-locked {
    background: rgba(255, 255, 255, 0.05);
    color: #5a5770;
}
```

**Step 4: Generate templ and verify it compiles**

Run: `task templ && go build ./...`
Expected: No errors

**Step 5: Commit**

```bash
git add internal/interface/web/templates/layout.templ internal/interface/web/server.go
git commit -m "feat: add Font Awesome CDN and gamification CSS"
```

---

## Task 12: Dashboard Template — Gamification Components

**Files:**
- Modify: `internal/interface/web/templates/dashboard.templ`

**Step 1: Add gamification components to the dashboard**

Insert these new sections at the top of the `Dashboard` templ function, before the Summary Cards, but after the `<div id="tab-content">`:

Replace the existing `<h2 class="title is-4">Dashboard</h2>` with:

```
<!-- Pool Health Score -->
<div class="box pv-neumorphic pv-health-score">
    <p class="pv-score-number">
        <span class={ healthScoreColor(data.HealthScore.Status) }>
            { fmt.Sprintf("%d", data.HealthScore.Score) }
        </span>
    </p>
    <p class="pv-score-label has-text-grey">{ data.HealthScore.Label }</p>
</div>
<!-- Streaks -->
if data.Streaks.TestingStreak > 1 || data.Streaks.TaskStreak > 1 {
    <div class="pv-streaks">
        if data.Streaks.TestingStreak > 1 {
            <span class="pv-streak-item">
                { fmt.Sprintf("%d week testing streak", data.Streaks.TestingStreak) }
            </span>
        }
        if data.Streaks.TaskStreak > 1 {
            <span class="pv-streak-item">
                { fmt.Sprintf("%d week task streak", data.Streaks.TaskStreak) }
            </span>
        }
    </div>
}
<!-- Milestones -->
<div class="pv-milestones mb-4">
    for _, m := range data.Milestones {
        @milestoneBadge(m)
    }
</div>
```

Add the `milestoneBadge` sub-component:

```
templ milestoneBadge(m templates.MilestoneBadge) {
    if m.Earned {
        <span class={ "pv-milestone-badge is-earned", templ.KV("is-new", m.IsNew) }>
            <i class={ m.Icon }></i>
            { m.Name }
        </span>
    } else {
        <span class="pv-milestone-badge is-locked">
            <i class={ m.Icon }></i>
            { m.Name }
        </span>
    }
}
```

Wait — `milestoneBadge` needs to be in the same `templates` package. Since `MilestoneBadge` is a type in `dashboard_types.go`, use it directly as `MilestoneBadge` (not `templates.MilestoneBadge`) inside the `.templ` file since it's in the same package.

Add the `healthScoreColor` helper to `internal/interface/web/templates/helpers.go`:

```go
func healthScoreColor(status string) string {
	switch status {
	case "danger":
		return "has-text-danger"
	case "warning":
		return "has-text-warning"
	default:
		return "has-text-success"
	}
}
```

Actually, this is the same as the existing `statusColor()` function. Just reuse `statusColor(data.HealthScore.Status)` directly.

**Step 2: Generate templates and verify compilation**

Run: `task templ && go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/interface/web/templates/dashboard.templ
git commit -m "feat: add health score, streaks, and milestone badges to dashboard"
```

---

## Task 13: Demo Seed — Milestones for Demo Users

**Files:**
- Modify: `internal/application/services/demo_seed_service.go`
- Modify: `internal/application/services/demo_cleanup_service.go`

**Step 1: Add milestone repo to DemoSeedService**

Add `milestoneRepo repositories.MilestoneRepository` field to the struct and constructor.

Update `NewDemoSeedService` to accept and store the milestone repo.

Add a `seedMilestones` method that earns a few milestones for the demo user (First Dip, Balanced, Stocked Up — since demo data satisfies these):

```go
func (s *DemoSeedService) seedMilestones(ctx context.Context, userID uuid.UUID) error {
	milestones := []entities.MilestoneKey{
		entities.MilestoneFirstDip,
		entities.MilestoneBalanced,
		entities.MilestoneStockedUp,
	}
	for _, key := range milestones {
		m := entities.NewMilestone(userID, key)
		if err := s.milestoneRepo.Create(ctx, m); err != nil {
			return err
		}
	}
	return nil
}
```

Call `s.seedMilestones(ctx, userID)` at the end of `Seed()`.

**Step 2: Add milestone cleanup to DemoCleanupService**

Add `milestoneRepo repositories.MilestoneRepository` field to `DemoCleanupService` struct and constructor.

In the `cleanup()` method, add before "Finally delete the user":
```go
_ = s.milestoneRepo.DeleteByUserID(ctx, user.ID)
```

**Step 3: Update wiring in `cmd/serve.go`**

Update `services.NewDemoSeedService(...)` call to pass `milestoneRepo`.
Update `services.NewDemoCleanupService(...)` call to pass `milestoneRepo`.

**Step 4: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 5: Commit**

```bash
git add internal/application/services/demo_seed_service.go internal/application/services/demo_cleanup_service.go cmd/serve.go
git commit -m "feat: seed milestones for demo users and clean up on expiry"
```

---

## Task 14: Build, Generate, and Manual Smoke Test

**Step 1: Generate templ files**

Run: `task templ`
Expected: No errors

**Step 2: Build**

Run: `task build`
Expected: No errors

**Step 3: Run all tests**

Run: `task test`
Expected: All tests pass

**Step 4: Manual smoke test with demo mode**

Run: `./bin/poolvibes serve --demo --db /tmp/test-gamification.db`

1. Open browser to `http://localhost:8080`
2. Sign up a new demo user
3. Verify dashboard shows:
   - Pool Health Score at top (should be high, e.g. 80-95 for demo data)
   - Testing streak visible (demo data has ~100 logs over 12 months)
   - Milestone badges visible (First Dip, Balanced, Stocked Up earned; others locked)
   - Existing summary cards, charts, and quick lists still work
4. Check dark mode appearance

Clean up: Kill the server, `rm /tmp/test-gamification.db`

**Step 5: Commit any generated `*_templ.go` files**

```bash
git add internal/interface/web/templates/*_templ.go
git commit -m "chore: regenerate templ output"
```

---

## Task 15: Update Documentation

**Files:**
- Modify: `docs/features/` — add gamification docs or update existing
- Modify: `docs/architecture.md` — mention gamification service
- Modify: `README.md` — mention gamification in features list

**Step 1: Add gamification to feature docs**

Create `docs/features/gamification.md` with a brief overview of Pool Health Score, Streaks, and Milestones.

**Step 2: Update architecture docs**

Add `GamificationService` to the services list in `docs/architecture.md` and mention the `user_milestones` table.

**Step 3: Update README**

Add "Pool Health Score & Milestones" to the features list in `README.md`.

**Step 4: Commit**

```bash
git add docs/ README.md
git commit -m "docs: add gamification feature documentation"
```
