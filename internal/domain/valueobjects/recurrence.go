package valueobjects

import (
	"fmt"
	"time"
)

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
)

type Recurrence struct {
	Frequency Frequency
	Interval  int
}

func NewRecurrence(frequency Frequency, interval int) (Recurrence, error) {
	if interval < 1 {
		return Recurrence{}, fmt.Errorf("interval must be at least 1")
	}
	switch frequency {
	case FrequencyDaily, FrequencyWeekly, FrequencyMonthly:
	default:
		return Recurrence{}, fmt.Errorf("invalid frequency: %s", frequency)
	}
	return Recurrence{Frequency: frequency, Interval: interval}, nil
}

func (r Recurrence) NextDueDate(from time.Time) time.Time {
	switch r.Frequency {
	case FrequencyDaily:
		return from.AddDate(0, 0, r.Interval)
	case FrequencyWeekly:
		return from.AddDate(0, 0, 7*r.Interval)
	case FrequencyMonthly:
		return from.AddDate(0, r.Interval, 0)
	default:
		return from
	}
}
