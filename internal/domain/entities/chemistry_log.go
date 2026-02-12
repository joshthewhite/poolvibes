package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChemistryLog struct {
	ID               uuid.UUID
	PH               float64
	FreeChlorine     float64
	CombinedChlorine float64
	TotalAlkalinity  float64
	CYA              float64
	CalciumHardness  float64
	Temperature      float64
	Notes            string
	TestedAt         time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewChemistryLog(ph, freeChlorine, combinedChlorine, totalAlkalinity, cya, calciumHardness, temperature float64, notes string, testedAt time.Time) *ChemistryLog {
	now := time.Now()
	return &ChemistryLog{
		ID:               uuid.Must(uuid.NewV7()),
		PH:               ph,
		FreeChlorine:     freeChlorine,
		CombinedChlorine: combinedChlorine,
		TotalAlkalinity:  totalAlkalinity,
		CYA:              cya,
		CalciumHardness:  calciumHardness,
		Temperature:      temperature,
		Notes:            notes,
		TestedAt:         testedAt,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (c *ChemistryLog) Validate() error {
	if c.PH < 0 || c.PH > 14 {
		return fmt.Errorf("pH must be between 0 and 14")
	}
	if c.FreeChlorine < 0 {
		return fmt.Errorf("free chlorine cannot be negative")
	}
	if c.CombinedChlorine < 0 {
		return fmt.Errorf("combined chlorine cannot be negative")
	}
	if c.TotalAlkalinity < 0 {
		return fmt.Errorf("total alkalinity cannot be negative")
	}
	if c.CYA < 0 {
		return fmt.Errorf("CYA cannot be negative")
	}
	if c.CalciumHardness < 0 {
		return fmt.Errorf("calcium hardness cannot be negative")
	}
	if c.TestedAt.IsZero() {
		return fmt.Errorf("tested_at is required")
	}
	return nil
}

func (c *ChemistryLog) PHInRange() bool { return c.PH >= 7.2 && c.PH <= 7.6 }
func (c *ChemistryLog) FreeChlorineInRange() bool {
	return c.FreeChlorine >= 1.0 && c.FreeChlorine <= 3.0
}
func (c *ChemistryLog) CombinedChlorineInRange() bool { return c.CombinedChlorine <= 0.5 }
func (c *ChemistryLog) TotalAlkalinityInRange() bool {
	return c.TotalAlkalinity >= 80 && c.TotalAlkalinity <= 120
}
func (c *ChemistryLog) CYAInRange() bool { return c.CYA >= 30 && c.CYA <= 50 }
func (c *ChemistryLog) CalciumHardnessInRange() bool {
	return c.CalciumHardness >= 200 && c.CalciumHardness <= 400
}
