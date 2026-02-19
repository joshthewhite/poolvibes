package handlers

import (
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
	"github.com/starfederation/datastar-go/datastar"
)

type DashboardHandler struct {
	chemSvc       *services.ChemistryService
	taskSvc       *services.TaskService
	chemicSvc     *services.ChemicalService
	milestoneRepo repositories.MilestoneRepository
}

func NewDashboardHandler(chemSvc *services.ChemistryService, taskSvc *services.TaskService, chemicSvc *services.ChemicalService, milestoneRepo repositories.MilestoneRepository) *DashboardHandler {
	return &DashboardHandler{chemSvc: chemSvc, taskSvc: taskSvc, chemicSvc: chemicSvc, milestoneRepo: milestoneRepo}
}

func (h *DashboardHandler) Page(w http.ResponseWriter, r *http.Request) {
	logs, _ := h.chemSvc.List(r.Context())
	tasks, _ := h.taskSvc.List(r.Context())
	chemicals, _ := h.chemicSvc.List(r.Context())

	data := buildDashboardData(logs, tasks, chemicals)

	// Gamification: health score, streaks, milestones
	now := time.Now()
	score := services.ComputeHealthScore(logs, tasks, chemicals, now)
	data.HealthScore = templates.HealthScoreSummary{
		Score:  score,
		Status: healthScoreStatus(score),
		Label:  healthScoreLabel(score),
	}
	data.Streaks = templates.StreaksSummary{
		TestingStreak: services.ComputeTestingStreak(logs, now),
		TaskStreak:    services.ComputeTaskStreak(tasks, now),
	}

	// User info for greeting
	user, err := services.UserFromContext(r.Context())
	if err == nil {
		data.Email = user.Email
	}

	// Milestones: check and persist newly earned
	if err == nil {
		existing, _ := h.milestoneRepo.FindAll(r.Context(), user.ID)
		earnedSet := make(map[entities.MilestoneKey]bool, len(existing))
		for _, m := range existing {
			earnedSet[m.Milestone] = true
		}

		newlyEarned := services.CheckMilestones(logs, tasks, chemicals, score, earnedSet)
		for _, key := range newlyEarned {
			m := entities.NewMilestone(user.ID, key)
			if err := h.milestoneRepo.Create(r.Context(), m); err != nil {
				slog.Error("Failed to persist milestone", "key", key, "error", err)
			}
			earnedSet[key] = true
		}

		newSet := make(map[entities.MilestoneKey]bool, len(newlyEarned))
		for _, key := range newlyEarned {
			newSet[key] = true
		}
		data.Milestones = buildMilestoneBadges(earnedSet, newSet)
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.Dashboard(data))
}

func buildDashboardData(logs []entities.ChemistryLog, tasks []entities.Task, chemicals []entities.Chemical) templates.DashboardData {
	data := templates.DashboardData{
		Chart: templates.ChartData{
			PHMin: 7.2,
			PHMax: 7.6,
			FCMin: 1.0,
			FCMax: 3.0,
		},
	}

	// Water quality & last tested
	if len(logs) > 0 {
		latest := logs[0] // logs are returned newest first
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

		status := "good"
		if inRange < 6 {
			status = "warning"
		}
		if inRange < 4 {
			status = "danger"
		}

		data.WaterQuality = templates.WaterQualitySummary{
			InRange: inRange,
			Total:   6,
			Status:  status,
			HasData: true,
		}

		// Last tested
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		tested := time.Date(latest.TestedAt.Year(), latest.TestedAt.Month(), latest.TestedAt.Day(), 0, 0, 0, 0, latest.TestedAt.Location())
		days := int(today.Sub(tested).Hours() / 24)

		var testedText string
		testedStatus := "good"
		switch {
		case days == 0:
			testedText = "Today"
		case days == 1:
			testedText = "Yesterday"
		default:
			testedText = fmt.Sprintf("%d days ago", days)
			if days > 3 {
				testedStatus = "warning"
			}
			if days > 7 {
				testedStatus = "danger"
			}
		}

		data.LastTested = templates.LastTestedSummary{
			Text:    testedText,
			Status:  testedStatus,
			HasData: true,
		}

		// Chart data — last 30 readings in chronological order
		chartLogs := make([]entities.ChemistryLog, len(logs))
		copy(chartLogs, logs)
		if len(chartLogs) > 30 {
			chartLogs = chartLogs[:30]
		}
		// Reverse to chronological order (logs come newest-first)
		for i, j := 0, len(chartLogs)-1; i < j; i, j = i+1, j-1 {
			chartLogs[i], chartLogs[j] = chartLogs[j], chartLogs[i]
		}

		labels := make([]string, len(chartLogs))
		phVals := make([]float64, len(chartLogs))
		fcVals := make([]float64, len(chartLogs))
		for i, l := range chartLogs {
			labels[i] = l.TestedAt.Format("Jan 2")
			phVals[i] = math.Round(l.PH*100) / 100
			fcVals[i] = math.Round(l.FreeChlorine*100) / 100
		}

		data.Chart.HasData = true
		data.Chart.SinglePoint = len(chartLogs) == 1
		data.Chart.Labels = labels
		data.Chart.PH = phVals
		data.Chart.FC = fcVals
	}

	// Tasks summary
	overdueCount := 0
	dueTodayCount := 0
	var upcomingTasks []entities.Task
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for i := range tasks {
		tasks[i].CheckOverdue()
		if tasks[i].Status == entities.TaskStatusCompleted {
			continue
		}
		if tasks[i].Status == entities.TaskStatusOverdue {
			overdueCount++
		}
		due := time.Date(tasks[i].DueDate.Year(), tasks[i].DueDate.Month(), tasks[i].DueDate.Day(), 0, 0, 0, 0, tasks[i].DueDate.Location())
		if today.Equal(due) {
			dueTodayCount++
		}
		upcomingTasks = append(upcomingTasks, tasks[i])
	}

	sort.Slice(upcomingTasks, func(i, j int) bool {
		return upcomingTasks[i].DueDate.Before(upcomingTasks[j].DueDate)
	})
	if len(upcomingTasks) > 7 {
		upcomingTasks = upcomingTasks[:7]
	}

	taskStatus := "good"
	if dueTodayCount > 0 {
		taskStatus = "warning"
	}
	if overdueCount > 0 {
		taskStatus = "danger"
	}

	data.Tasks = templates.TaskSummary{
		OverdueCount:  overdueCount,
		DueTodayCount: dueTodayCount,
		Status:        taskStatus,
		HasData:       len(tasks) > 0,
	}
	data.UpcomingTasks = upcomingTasks

	// Low stock
	var lowStockChemicals []entities.Chemical
	for _, c := range chemicals {
		if c.IsLowStock() {
			lowStockChemicals = append(lowStockChemicals, c)
		}
	}

	lowStockStatus := "good"
	if len(lowStockChemicals) > 0 {
		lowStockStatus = "warning"
	}
	if len(lowStockChemicals) > 2 {
		lowStockStatus = "danger"
	}

	data.LowStock = templates.LowStockSummary{
		Count:   len(lowStockChemicals),
		Status:  lowStockStatus,
		HasData: len(chemicals) > 0,
	}
	data.LowStockChemicals = lowStockChemicals

	return data
}

func healthScoreStatus(score int) string {
	switch {
	case score >= 80:
		return "good"
	case score >= 50:
		return "warning"
	default:
		return "danger"
	}
}

func healthScoreLabel(score int) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 80:
		return "Great"
	case score >= 60:
		return "Needs attention"
	case score >= 40:
		return "Falling behind"
	default:
		return "Critical"
	}
}

var milestoneInfo = map[entities.MilestoneKey]struct {
	Name string
	Icon string
}{
	entities.MilestoneFirstDip:    {"First Dip", "fa-solid fa-droplet"},
	entities.MilestoneBalanced:    {"Balanced", "fa-solid fa-scale-balanced"},
	entities.MilestoneConsistent:  {"Consistent", "fa-solid fa-calendar-check"},
	entities.MilestoneDevoted:     {"Devoted", "fa-solid fa-award"},
	entities.MilestoneOnIt:        {"On It", "fa-solid fa-clipboard-check"},
	entities.MilestoneStockedUp:   {"Stocked Up", "fa-solid fa-boxes-stacked"},
	entities.MilestoneCleanRecord: {"Clean Record", "fa-solid fa-shield-halved"},
	entities.MilestonePoolPro:     {"Pool Pro", "fa-solid fa-trophy"},
}

func buildMilestoneBadges(earned, newlyEarned map[entities.MilestoneKey]bool) []templates.MilestoneBadge {
	var badges []templates.MilestoneBadge
	for _, key := range entities.AllMilestones() {
		info := milestoneInfo[key]
		badges = append(badges, templates.MilestoneBadge{
			Key:    string(key),
			Name:   info.Name,
			Icon:   info.Icon,
			Earned: earned[key],
			IsNew:  newlyEarned[key],
		})
	}
	return badges
}
