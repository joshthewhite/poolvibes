package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func makeLog(ph, fc, cc, ta, cya, ch float64) *ChemistryLog {
	return &ChemistryLog{
		ID:               uuid.Must(uuid.NewV7()),
		UserID:           uuid.Must(uuid.NewV7()),
		PH:               ph,
		FreeChlorine:     fc,
		CombinedChlorine: cc,
		TotalAlkalinity:  ta,
		CYA:              cya,
		CalciumHardness:  ch,
		Temperature:      82,
		TestedAt:         time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func TestGenerateTreatmentPlan_AllInRange(t *testing.T) {
	log := makeLog(7.4, 2.0, 0.2, 100, 40, 300)
	plan := GenerateTreatmentPlan(log, 15000)
	if len(plan.Steps) != 0 {
		t.Errorf("expected 0 steps for in-range values, got %d", len(plan.Steps))
		for _, s := range plan.Steps {
			t.Logf("  step: %s", s.Problem)
		}
	}
}

func TestGenerateTreatmentPlan(t *testing.T) {
	tests := []struct {
		name        string
		ph          float64
		fc          float64
		cc          float64
		ta          float64
		cya         float64
		ch          float64
		poolGallons int
		wantProblem string
		wantChem    string
	}{
		{
			name:        "high pH",
			ph:          8.0,
			fc:          2.0,
			cc:          0.2,
			ta:          100,
			cya:         40,
			ch:          300,
			poolGallons: 10000,
			wantProblem: "High pH",
			wantChem:    "Muriatic acid",
		},
		{
			name:        "low pH",
			ph:          6.8,
			fc:          2.0,
			cc:          0.2,
			ta:          100,
			cya:         40,
			ch:          300,
			poolGallons: 10000,
			wantProblem: "Low pH",
			wantChem:    "Soda ash",
		},
		{
			name:        "low free chlorine",
			ph:          7.4,
			fc:          0.5,
			cc:          0.2,
			ta:          100,
			cya:         40,
			ch:          300,
			poolGallons: 10000,
			wantProblem: "Low free chlorine",
			wantChem:    "Calcium hypochlorite",
		},
		{
			name:        "high combined chlorine",
			ph:          7.4,
			fc:          2.0,
			cc:          1.0,
			ta:          100,
			cya:         40,
			ch:          300,
			poolGallons: 10000,
			wantProblem: "High combined chlorine",
			wantChem:    "Calcium hypochlorite",
		},
		{
			name:        "low total alkalinity",
			ph:          7.4,
			fc:          2.0,
			cc:          0.2,
			ta:          60,
			cya:         40,
			ch:          300,
			poolGallons: 10000,
			wantProblem: "Low total alkalinity",
			wantChem:    "Baking soda",
		},
		{
			name:        "high total alkalinity",
			ph:          7.4,
			fc:          2.0,
			cc:          0.2,
			ta:          160,
			cya:         40,
			ch:          300,
			poolGallons: 10000,
			wantProblem: "High total alkalinity",
			wantChem:    "Muriatic acid",
		},
		{
			name:        "low CYA",
			ph:          7.4,
			fc:          2.0,
			cc:          0.2,
			ta:          100,
			cya:         10,
			ch:          300,
			poolGallons: 10000,
			wantProblem: "Low CYA (stabilizer)",
			wantChem:    "Cyanuric acid",
		},
		{
			name:        "low calcium hardness",
			ph:          7.4,
			fc:          2.0,
			cc:          0.2,
			ta:          100,
			cya:         40,
			ch:          100,
			poolGallons: 10000,
			wantProblem: "Low calcium hardness",
			wantChem:    "Calcium chloride",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := makeLog(tt.ph, tt.fc, tt.cc, tt.ta, tt.cya, tt.ch)
			plan := GenerateTreatmentPlan(log, tt.poolGallons)

			found := false
			for _, step := range plan.Steps {
				if step.Problem == tt.wantProblem {
					found = true
					if !strings.Contains(step.Chemical, tt.wantChem) {
						t.Errorf("expected chemical containing %q, got %q", tt.wantChem, step.Chemical)
					}
					if step.Amount == "" {
						t.Error("expected non-empty amount")
					}
					if step.Instructions == "" {
						t.Error("expected non-empty instructions")
					}
				}
			}
			if !found {
				t.Errorf("expected step with problem %q, got steps: %v", tt.wantProblem, stepNames(plan.Steps))
			}
		})
	}
}

func TestGenerateTreatmentPlan_ScalesWithPoolSize(t *testing.T) {
	log := makeLog(8.0, 2.0, 0.2, 100, 40, 300) // high pH only

	plan10k := GenerateTreatmentPlan(log, 10000)
	plan20k := GenerateTreatmentPlan(log, 20000)

	if len(plan10k.Steps) != 1 || len(plan20k.Steps) != 1 {
		t.Fatal("expected exactly 1 step each")
	}

	// The amounts should be different (20k should use more chemical)
	if plan10k.Steps[0].Amount == plan20k.Steps[0].Amount {
		t.Errorf("expected different amounts for different pool sizes, both got %s", plan10k.Steps[0].Amount)
	}
}

func TestGenerateTreatmentPlan_MultipleIssues(t *testing.T) {
	// Everything out of range
	log := makeLog(8.2, 0.3, 1.5, 50, 10, 100)
	plan := GenerateTreatmentPlan(log, 15000)

	if len(plan.Steps) < 5 {
		t.Errorf("expected at least 5 steps for multiple issues, got %d", len(plan.Steps))
	}
}

func stepNames(steps []TreatmentStep) []string {
	names := make([]string, len(steps))
	for i, s := range steps {
		names[i] = s.Problem
	}
	return names
}
