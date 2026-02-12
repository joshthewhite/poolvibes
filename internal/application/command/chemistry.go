package command

import "time"

type CreateChemistryLog struct {
	PH               float64
	FreeChlorine     float64
	CombinedChlorine float64
	TotalAlkalinity  float64
	CYA              float64
	CalciumHardness  float64
	Temperature      float64
	Notes            string
	TestedAt         time.Time
}

type UpdateChemistryLog struct {
	ID               string
	PH               float64
	FreeChlorine     float64
	CombinedChlorine float64
	TotalAlkalinity  float64
	CYA              float64
	CalciumHardness  float64
	Temperature      float64
	Notes            string
	TestedAt         time.Time
}
