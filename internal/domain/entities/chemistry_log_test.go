package entities

import "testing"

func TestChemistryLog_PHInRange(t *testing.T) {
	tests := []struct {
		name string
		ph   float64
		want bool
	}{
		{"below range", 7.0, false},
		{"at lower bound", 7.2, true},
		{"mid range", 7.4, true},
		{"at upper bound", 7.6, true},
		{"above range", 7.8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChemistryLog{PH: tt.ph}
			if got := c.PHInRange(); got != tt.want {
				t.Errorf("PHInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChemistryLog_FreeChlorineInRange(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want bool
	}{
		{"below range", 0.5, false},
		{"at lower bound", 1.0, true},
		{"mid range", 2.0, true},
		{"at upper bound", 3.0, true},
		{"above range", 3.5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChemistryLog{FreeChlorine: tt.val}
			if got := c.FreeChlorineInRange(); got != tt.want {
				t.Errorf("FreeChlorineInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChemistryLog_CombinedChlorineInRange(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want bool
	}{
		{"zero", 0, true},
		{"at upper bound", 0.5, true},
		{"above range", 0.6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChemistryLog{CombinedChlorine: tt.val}
			if got := c.CombinedChlorineInRange(); got != tt.want {
				t.Errorf("CombinedChlorineInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChemistryLog_TotalAlkalinityInRange(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want bool
	}{
		{"below range", 70, false},
		{"at lower bound", 80, true},
		{"mid range", 100, true},
		{"at upper bound", 120, true},
		{"above range", 130, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChemistryLog{TotalAlkalinity: tt.val}
			if got := c.TotalAlkalinityInRange(); got != tt.want {
				t.Errorf("TotalAlkalinityInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChemistryLog_CYAInRange(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want bool
	}{
		{"below range", 20, false},
		{"at lower bound", 30, true},
		{"mid range", 40, true},
		{"at upper bound", 50, true},
		{"above range", 60, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChemistryLog{CYA: tt.val}
			if got := c.CYAInRange(); got != tt.want {
				t.Errorf("CYAInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChemistryLog_CalciumHardnessInRange(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want bool
	}{
		{"below range", 150, false},
		{"at lower bound", 200, true},
		{"mid range", 300, true},
		{"at upper bound", 400, true},
		{"above range", 450, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChemistryLog{CalciumHardness: tt.val}
			if got := c.CalciumHardnessInRange(); got != tt.want {
				t.Errorf("CalciumHardnessInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
