package templates

import (
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	// Fixed reference time: Feb 15, 2026 14:00:00 UTC
	now := time.Date(2026, 2, 15, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{"just now", now.Add(-30 * time.Second), "just now"},
		{"1 minute ago", now.Add(-1 * time.Minute), "1m ago"},
		{"5 minutes ago", now.Add(-5 * time.Minute), "5m ago"},
		{"59 minutes ago", now.Add(-59 * time.Minute), "59m ago"},
		{"1 hour ago", now.Add(-1 * time.Hour), "1h ago"},
		{"3 hours ago", now.Add(-3 * time.Hour), "3h ago"},
		{"23 hours ago", now.Add(-23 * time.Hour), "23h ago"},
		{"yesterday", time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC), "yesterday"},
		{"2 days ago", time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC), "2 days ago"},
		{"6 days ago", time.Date(2026, 2, 9, 10, 0, 0, 0, time.UTC), "6 days ago"},
		{"1 week ago", time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC), "1 week ago"},
		{"13 days ago", time.Date(2026, 2, 2, 10, 0, 0, 0, time.UTC), "1 week ago"},
		{"same year older", time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC), "Jan 10"},
		{"different year", time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC), "Jun 15, 2025"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := relativeTimeFrom(tt.t, now); got != tt.want {
				t.Errorf("relativeTimeFrom() = %q, want %q", got, tt.want)
			}
		})
	}
}
