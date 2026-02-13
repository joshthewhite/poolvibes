package entities

import (
	"testing"
	"time"
)

func TestEquipment_IsWarrantyActive(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	past := time.Now().Add(-24 * time.Hour)

	tests := []struct {
		name   string
		expiry *time.Time
		want   bool
	}{
		{"nil expiry", nil, false},
		{"future expiry", &future, true},
		{"past expiry", &past, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Equipment{WarrantyExpiry: tt.expiry}
			if got := e.IsWarrantyActive(); got != tt.want {
				t.Errorf("IsWarrantyActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
