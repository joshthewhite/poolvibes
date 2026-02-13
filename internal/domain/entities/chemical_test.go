package entities

import (
	"testing"

	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

func TestChemical_Validate(t *testing.T) {
	tests := []struct {
		name    string
		chem    Chemical
		wantErr string
	}{
		{
			name:    "missing name",
			chem:    Chemical{Type: ChemicalTypeSanitizer},
			wantErr: "name is required",
		},
		{
			name:    "invalid type",
			chem:    Chemical{Name: "Chlorine", Type: ChemicalType("invalid")},
			wantErr: "invalid chemical type: invalid",
		},
		{
			name:    "negative threshold",
			chem:    Chemical{Name: "Chlorine", Type: ChemicalTypeSanitizer, AlertThreshold: -1},
			wantErr: "alert threshold cannot be negative",
		},
		{
			name: "valid",
			chem: Chemical{Name: "Chlorine", Type: ChemicalTypeSanitizer, AlertThreshold: 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chem.Validate()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantErr {
					t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestChemical_IsLowStock(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		threshold float64
		want      bool
	}{
		{"below threshold", 3, 5, true},
		{"at threshold", 5, 5, true},
		{"above threshold", 7, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chemical{
				Stock:          valueobjects.Quantity{Amount: tt.amount, Unit: valueobjects.UnitPounds},
				AlertThreshold: tt.threshold,
			}
			if got := c.IsLowStock(); got != tt.want {
				t.Errorf("IsLowStock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChemical_AdjustStock(t *testing.T) {
	tests := []struct {
		name        string
		startStock  float64
		delta       float64
		wantStock   float64
		wantErr     bool
	}{
		{"positive delta", 10, 5, 15, false},
		{"negative delta", 10, -3, 7, false},
		{"delta to zero", 5, -5, 0, false},
		{"delta below zero", 5, -6, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chemical{
				Stock: valueobjects.Quantity{Amount: tt.startStock, Unit: valueobjects.UnitPounds},
			}
			err := c.AdjustStock(tt.delta)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.Stock.Amount != tt.wantStock {
				t.Errorf("Stock.Amount = %v, want %v", c.Stock.Amount, tt.wantStock)
			}
		})
	}
}
