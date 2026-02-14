package entities

import (
	"fmt"
	"math"
)

type TreatmentStep struct {
	Problem      string
	Explanation  string
	Chemical     string
	Amount       string
	MaxDose      string
	Instructions string
}

type TreatmentPlan struct {
	Steps       []TreatmentStep
	PoolGallons int
}

// GenerateTreatmentPlan computes chemical dosages to correct out-of-range
// readings. All dosage formulas are per 10,000 gallons, scaled to pool volume.
func GenerateTreatmentPlan(log *ChemistryLog, poolGallons int) *TreatmentPlan {
	plan := &TreatmentPlan{PoolGallons: poolGallons}
	scale := float64(poolGallons) / 10000.0

	// High pH (>7.6) → muriatic acid (31.45% HCl)
	// ~26 fl oz per 10k gal lowers pH by 0.1
	if log.PH > 7.6 {
		drop := log.PH - 7.4 // target 7.4
		ozPer10k := drop / 0.1 * 26.0
		totalOz := ozPer10k * scale
		plan.Steps = append(plan.Steps, TreatmentStep{
			Problem:      "High pH",
			Explanation:  "High pH reduces chlorine effectiveness, causes cloudy water, and promotes scale formation.",
			Chemical:     "Muriatic acid (31.45% HCl)",
			Amount:       fmtFlOz(totalOz),
			MaxDose:      fmtFlOz(math.Min(totalOz, 32*scale)),
			Instructions: "With pump running, pour slowly into the deep end away from walls and fittings. Wait 4 hours and retest before adding more.",
		})
	}

	// Low pH (<7.2) → soda ash (sodium carbonate)
	// ~6 oz (weight) per 10k gal raises pH by 0.1
	if log.PH < 7.2 {
		raise := 7.4 - log.PH // target 7.4
		ozPer10k := raise / 0.1 * 6.0
		totalOz := ozPer10k * scale
		plan.Steps = append(plan.Steps, TreatmentStep{
			Problem:      "Low pH",
			Explanation:  "Low pH causes eye/skin irritation, corrodes equipment, and etches plaster surfaces.",
			Chemical:     "Soda ash (sodium carbonate)",
			Amount:       fmtWeightOz(totalOz),
			MaxDose:      fmtWeightOz(math.Min(totalOz, 16*scale)),
			Instructions: "Pre-dissolve in a bucket of pool water. Pour around the pool perimeter with pump running. Wait 4 hours and retest.",
		})
	}

	// Low free chlorine (<1.0 ppm) → calcium hypochlorite (cal-hypo 73%)
	// ~2 oz (weight) per 10k gal raises FC by 1 ppm
	if log.FreeChlorine < 1.0 {
		raise := 2.0 - log.FreeChlorine // target 2.0 ppm
		ozPer10k := raise * 2.0
		totalOz := ozPer10k * scale
		plan.Steps = append(plan.Steps, TreatmentStep{
			Problem:      "Low free chlorine",
			Explanation:  "Insufficient free chlorine allows algae and bacteria to grow, making the pool unsafe for swimming.",
			Chemical:     "Calcium hypochlorite (cal-hypo 73%)",
			Amount:       fmtWeightOz(totalOz),
			MaxDose:      fmtWeightOz(math.Min(totalOz, 4*scale)),
			Instructions: "Pre-dissolve in a bucket of water. Pour around the pool perimeter with pump running. Do not swim for at least 30 minutes or until FC drops below 5 ppm.",
		})
	}

	// High combined chlorine (>0.5 ppm) → breakpoint chlorination (shock)
	// Need to raise FC to 10x the CC level. ~2 oz cal-hypo per 10k gal per 1 ppm
	if log.CombinedChlorine > 0.5 {
		targetFC := log.CombinedChlorine * 10
		raise := targetFC - log.FreeChlorine
		if raise < 0 {
			raise = 0
		}
		ozPer10k := raise * 2.0
		totalOz := ozPer10k * scale
		if totalOz > 0 {
			plan.Steps = append(plan.Steps, TreatmentStep{
				Problem:      "High combined chlorine",
				Explanation:  "Combined chlorine (chloramines) causes the harsh chlorine smell and eye irritation. Breakpoint chlorination destroys chloramines.",
				Chemical:     "Calcium hypochlorite (cal-hypo 73%)",
				Amount:       fmtWeightOz(totalOz),
				MaxDose:      fmtWeightOz(math.Min(totalOz, 8*scale)),
				Instructions: "This is a shock treatment. Pre-dissolve in a bucket and distribute around the pool at dusk. Run pump overnight. Do not swim until FC drops below 5 ppm.",
			})
		}
	}

	// Low total alkalinity (<80 ppm) → baking soda (sodium bicarbonate)
	// ~1.5 lbs per 10k gal raises TA by 10 ppm
	if log.TotalAlkalinity < 80 {
		raise := 100 - log.TotalAlkalinity // target 100 ppm
		lbsPer10k := raise / 10.0 * 1.5
		totalLbs := lbsPer10k * scale
		plan.Steps = append(plan.Steps, TreatmentStep{
			Problem:      "Low total alkalinity",
			Explanation:  "Low alkalinity causes pH to fluctuate rapidly, leading to corrosion and difficulty maintaining balance.",
			Chemical:     "Baking soda (sodium bicarbonate)",
			Amount:       fmtLbs(totalLbs),
			MaxDose:      fmtLbs(math.Min(totalLbs, 3*scale)),
			Instructions: "Broadcast over the pool surface with pump running. Add no more than 3 lbs per 10,000 gallons at a time. Wait 6 hours and retest before adding more.",
		})
	}

	// High total alkalinity (>120 ppm) → muriatic acid
	// ~26 fl oz per 10k gal lowers TA by ~10 ppm (also lowers pH)
	if log.TotalAlkalinity > 120 {
		drop := log.TotalAlkalinity - 100 // target 100 ppm
		ozPer10k := drop / 10.0 * 26.0
		totalOz := ozPer10k * scale
		plan.Steps = append(plan.Steps, TreatmentStep{
			Problem:      "High total alkalinity",
			Explanation:  "High alkalinity makes it difficult to adjust pH and can cause cloudy water and scale buildup.",
			Chemical:     "Muriatic acid (31.45% HCl)",
			Amount:       fmtFlOz(totalOz),
			MaxDose:      fmtFlOz(math.Min(totalOz, 32*scale)),
			Instructions: "Pour slowly in one spot in the deep end with pump off, then turn pump on after 1 hour. This technique helps lower TA without dropping pH as much. Wait 6 hours and retest.",
		})
	}

	// Low CYA (<30 ppm) → cyanuric acid (stabilizer)
	// ~13 oz (weight) per 10k gal raises CYA by 10 ppm
	if log.CYA < 30 {
		raise := 40 - log.CYA // target 40 ppm
		ozPer10k := raise / 10.0 * 13.0
		totalOz := ozPer10k * scale
		plan.Steps = append(plan.Steps, TreatmentStep{
			Problem:      "Low CYA (stabilizer)",
			Explanation:  "Without adequate CYA, sunlight rapidly destroys chlorine. Your pool can lose most of its chlorine in just a few hours.",
			Chemical:     "Cyanuric acid (stabilizer)",
			Amount:       fmtWeightOz(totalOz),
			MaxDose:      fmtWeightOz(math.Min(totalOz, 16*scale)),
			Instructions: "Place in a sock or mesh bag in front of a return jet, or add to the skimmer basket. CYA dissolves slowly — allow 48 hours to fully dissolve and circulate before retesting.",
		})
	}

	// Low calcium hardness (<200 ppm) → calcium chloride
	// ~1.25 lbs per 10k gal raises CH by 10 ppm
	if log.CalciumHardness < 200 {
		raise := 300 - log.CalciumHardness // target 300 ppm
		lbsPer10k := raise / 10.0 * 1.25
		totalLbs := lbsPer10k * scale
		plan.Steps = append(plan.Steps, TreatmentStep{
			Problem:      "Low calcium hardness",
			Explanation:  "Low calcium causes the water to become aggressive, dissolving calcium from plaster, grout, and equipment.",
			Chemical:     "Calcium chloride",
			Amount:       fmtLbs(totalLbs),
			MaxDose:      fmtLbs(math.Min(totalLbs, 2.5*scale)),
			Instructions: "Pre-dissolve in a bucket of pool water (it generates heat — use caution). Pour around the pool perimeter with pump running. Add no more than 2.5 lbs per 10,000 gallons at a time. Wait 6 hours and retest.",
		})
	}

	return plan
}

func fmtFlOz(oz float64) string {
	if oz >= 128 {
		gal := oz / 128
		return fmt.Sprintf("%.1f gal", gal)
	}
	return fmt.Sprintf("%.0f fl oz", math.Round(oz))
}

func fmtWeightOz(oz float64) string {
	if oz >= 16 {
		lbs := oz / 16
		return fmt.Sprintf("%.1f lbs", lbs)
	}
	return fmt.Sprintf("%.0f oz", math.Round(oz))
}

func fmtLbs(lbs float64) string {
	if lbs < 1 {
		oz := lbs * 16
		return fmt.Sprintf("%.0f oz", math.Round(oz))
	}
	return fmt.Sprintf("%.1f lbs", lbs)
}
