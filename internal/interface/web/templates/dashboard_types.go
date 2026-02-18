package templates

import "github.com/joshthewhite/poolvibes/internal/domain/entities"

type DashboardData struct {
	WaterQuality      WaterQualitySummary
	LastTested        LastTestedSummary
	Tasks             TaskSummary
	LowStock          LowStockSummary
	Chart             ChartData
	UpcomingTasks     []entities.Task
	LowStockChemicals []entities.Chemical
	HealthScore       HealthScoreSummary
	Streaks           StreaksSummary
	Milestones        []MilestoneBadge
}

type WaterQualitySummary struct {
	InRange int
	Total   int
	Status  string
	HasData bool
}

type LastTestedSummary struct {
	Text    string
	Status  string
	HasData bool
}

type TaskSummary struct {
	OverdueCount  int
	DueTodayCount int
	Status        string
	HasData       bool
}

type LowStockSummary struct {
	Count   int
	Status  string
	HasData bool
}

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
	Key    string
	Name   string
	Icon   string // Font Awesome class
	Earned bool
	IsNew  bool // just earned this load
}

type ChartData struct {
	HasData     bool      `json:"hasData"`
	SinglePoint bool      `json:"singlePoint"`
	Labels      []string  `json:"labels"`
	PH          []float64 `json:"ph"`
	FC          []float64 `json:"fc"`
	PHMin       float64   `json:"phMin"`
	PHMax       float64   `json:"phMax"`
	FCMin       float64   `json:"fcMin"`
	FCMax       float64   `json:"fcMax"`
}
