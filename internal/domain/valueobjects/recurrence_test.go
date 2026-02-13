package valueobjects

import (
	"testing"
	"time"
)

func TestNewRecurrence(t *testing.T) {
	tests := []struct {
		name      string
		frequency Frequency
		interval  int
		wantErr   bool
	}{
		{"daily", FrequencyDaily, 1, false},
		{"weekly", FrequencyWeekly, 2, false},
		{"monthly", FrequencyMonthly, 3, false},
		{"invalid frequency", Frequency("yearly"), 1, true},
		{"zero interval", FrequencyDaily, 0, true},
		{"negative interval", FrequencyDaily, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRecurrence(tt.frequency, tt.interval)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.Frequency != tt.frequency {
				t.Errorf("Frequency = %v, want %v", r.Frequency, tt.frequency)
			}
			if r.Interval != tt.interval {
				t.Errorf("Interval = %v, want %v", r.Interval, tt.interval)
			}
		})
	}
}

func TestRecurrence_NextDueDate(t *testing.T) {
	base := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		frequency Frequency
		interval  int
		want      time.Time
	}{
		{"daily interval 1", FrequencyDaily, 1, time.Date(2025, 1, 16, 12, 0, 0, 0, time.UTC)},
		{"daily interval 3", FrequencyDaily, 3, time.Date(2025, 1, 18, 12, 0, 0, 0, time.UTC)},
		{"weekly interval 1", FrequencyWeekly, 1, time.Date(2025, 1, 22, 12, 0, 0, 0, time.UTC)},
		{"weekly interval 2", FrequencyWeekly, 2, time.Date(2025, 1, 29, 12, 0, 0, 0, time.UTC)},
		{"monthly interval 1", FrequencyMonthly, 1, time.Date(2025, 2, 15, 12, 0, 0, 0, time.UTC)},
		{"monthly interval 3", FrequencyMonthly, 3, time.Date(2025, 4, 15, 12, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Recurrence{Frequency: tt.frequency, Interval: tt.interval}
			got := r.NextDueDate(base)
			if !got.Equal(tt.want) {
				t.Errorf("NextDueDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
