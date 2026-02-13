package valueobjects

import "testing"

func TestNewQuantity(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		unit    Unit
		wantErr bool
	}{
		{"valid pounds", 10.5, UnitPounds, false},
		{"valid ounces", 0, UnitOunces, false},
		{"valid gallons", 5.0, UnitGallons, false},
		{"valid liters", 1.0, UnitLiters, false},
		{"valid kg", 2.5, UnitKg, false},
		{"negative amount", -1.0, UnitPounds, true},
		{"invalid unit", 1.0, Unit("cups"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuantity(tt.amount, tt.unit)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if q.Amount != tt.amount {
				t.Errorf("Amount = %v, want %v", q.Amount, tt.amount)
			}
			if q.Unit != tt.unit {
				t.Errorf("Unit = %v, want %v", q.Unit, tt.unit)
			}
		})
	}
}
