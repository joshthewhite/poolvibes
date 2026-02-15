package services

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/domain/repositories"
	"github.com/joshthewhite/poolvibes/internal/domain/valueobjects"
)

type DemoSeedService struct {
	chemLogRepo repositories.ChemistryLogRepository
	taskRepo    repositories.TaskRepository
	equipRepo   repositories.EquipmentRepository
	srRepo      repositories.ServiceRecordRepository
	chemRepo    repositories.ChemicalRepository
}

func NewDemoSeedService(
	chemLogRepo repositories.ChemistryLogRepository,
	taskRepo repositories.TaskRepository,
	equipRepo repositories.EquipmentRepository,
	srRepo repositories.ServiceRecordRepository,
	chemRepo repositories.ChemicalRepository,
) *DemoSeedService {
	return &DemoSeedService{
		chemLogRepo: chemLogRepo,
		taskRepo:    taskRepo,
		equipRepo:   equipRepo,
		srRepo:      srRepo,
		chemRepo:    chemRepo,
	}
}

func (s *DemoSeedService) Seed(ctx context.Context, userID uuid.UUID) error {
	if err := s.seedChemistryLogs(ctx, userID); err != nil {
		return fmt.Errorf("seeding chemistry logs: %w", err)
	}
	if err := s.seedTasks(ctx, userID); err != nil {
		return fmt.Errorf("seeding tasks: %w", err)
	}
	equipIDs, err := s.seedEquipment(ctx, userID)
	if err != nil {
		return fmt.Errorf("seeding equipment: %w", err)
	}
	if err := s.seedServiceRecords(ctx, userID, equipIDs); err != nil {
		return fmt.Errorf("seeding service records: %w", err)
	}
	if err := s.seedChemicals(ctx, userID); err != nil {
		return fmt.Errorf("seeding chemicals: %w", err)
	}
	return nil
}

func (s *DemoSeedService) seedChemistryLogs(ctx context.Context, userID uuid.UUID) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	now := time.Now()

	for i := 0; i < 100; i++ {
		// Spread over ~12 months, roughly 2x per week
		daysAgo := int(float64(i) * 3.65)
		testedAt := now.AddDate(0, 0, -daysAgo)

		// Seasonal temperature variation (warmer in summer, cooler in winter)
		month := testedAt.Month()
		baseTemp := 75.0
		switch {
		case month >= time.June && month <= time.August:
			baseTemp = 85.0
		case month == time.May || month == time.September:
			baseTemp = 78.0
		case month >= time.November || month <= time.February:
			baseTemp = 62.0
		}

		// Occasionally produce out-of-range readings (~10% of the time)
		outOfRange := rng.Float64() < 0.10

		ph := clampF(randNorm(rng, 7.4, 0.15), 6.8, 8.0)
		freeChlorine := clampF(randNorm(rng, 2.0, 0.5), 0.5, 5.0)
		combinedChlorine := clampF(randNorm(rng, 0.2, 0.1), 0.0, 1.0)
		alkalinity := clampF(randNorm(rng, 100, 12), 60, 140)
		cya := clampF(randNorm(rng, 40, 8), 15, 70)
		hardness := clampF(randNorm(rng, 300, 50), 100, 500)
		temp := clampF(randNorm(rng, baseTemp, 3), 55, 95)

		if outOfRange {
			// Push one value out of ideal range
			switch rng.Intn(3) {
			case 0:
				ph = clampF(randNorm(rng, 7.8, 0.2), 7.6, 8.2)
			case 1:
				freeChlorine = clampF(randNorm(rng, 0.5, 0.3), 0.1, 0.9)
			case 2:
				alkalinity = clampF(randNorm(rng, 65, 5), 50, 75)
			}
		}

		log := entities.NewChemistryLog(
			userID, ph, freeChlorine, combinedChlorine,
			alkalinity, cya, hardness, temp, "", testedAt,
		)
		if err := s.chemLogRepo.Create(ctx, log); err != nil {
			return err
		}
	}
	return nil
}

func (s *DemoSeedService) seedTasks(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	tasks := []struct {
		name       string
		desc       string
		freq       valueobjects.Frequency
		interval   int
		dueDaysOut int
	}{
		{"Test water chemistry", "Test pH, chlorine, alkalinity, and CYA levels", valueobjects.FrequencyWeekly, 1, 2},
		{"Clean skimmer baskets", "Remove debris from all skimmer baskets", valueobjects.FrequencyWeekly, 1, -1},
		{"Backwash filter", "Backwash when pressure rises 8-10 PSI above clean", valueobjects.FrequencyMonthly, 1, 12},
		{"Check pump pressure", "Record pump PSI and check for unusual readings", valueobjects.FrequencyWeekly, 1, 4},
		{"Brush pool walls", "Brush walls and tile line to prevent algae buildup", valueobjects.FrequencyWeekly, 1, -3},
		{"Inspect equipment", "Visual inspection of pump, filter, heater, and chlorinator", valueobjects.FrequencyMonthly, 1, 20},
	}

	for _, t := range tasks {
		rec, _ := valueobjects.NewRecurrence(t.freq, t.interval)
		dueDate := now.AddDate(0, 0, t.dueDaysOut)
		task := entities.NewTask(userID, t.name, t.desc, rec, dueDate)
		if t.dueDaysOut < 0 {
			task.Status = entities.TaskStatusOverdue
		}
		if err := s.taskRepo.Create(ctx, task); err != nil {
			return err
		}
	}
	return nil
}

func (s *DemoSeedService) seedEquipment(ctx context.Context, userID uuid.UUID) (map[string]uuid.UUID, error) {
	now := time.Now()
	ids := make(map[string]uuid.UUID)

	items := []struct {
		name         string
		category     entities.EquipmentCategory
		manufacturer string
		model        string
		serial       string
		installYears int
		warrantyLeft int // months from now, 0 means no warranty
	}{
		{"Variable Speed Pump", entities.CategoryPump, "Pentair", "IntelliFlo VSF 011056", "VSF-2024-38291", 2, 12},
		{"Sand Filter", entities.CategoryFilter, "Hayward", "Pro Series S244T", "HW-S244-77201", 3, 0},
		{"Salt Chlorinator", entities.CategoryChlorinator, "Hayward", "AquaRite 925", "AQR-925-55103", 1, 24},
		{"Robotic Cleaner", entities.CategoryCleaner, "Dolphin", "Nautilus CC Plus", "DLP-NCC-40822", 1, 6},
		{"Gas Heater", entities.CategoryHeater, "Raypak", "336A Digital", "RP-336-19455", 4, 0},
	}

	for _, item := range items {
		installDate := now.AddDate(-item.installYears, 0, 0)
		var warrantyExpiry *time.Time
		if item.warrantyLeft > 0 {
			exp := now.AddDate(0, item.warrantyLeft, 0)
			warrantyExpiry = &exp
		}
		equip := entities.NewEquipment(
			userID, item.name, item.category,
			item.manufacturer, item.model, item.serial,
			&installDate, warrantyExpiry,
		)
		if err := s.equipRepo.Create(ctx, equip); err != nil {
			return nil, err
		}
		ids[item.name] = equip.ID
	}
	return ids, nil
}

func (s *DemoSeedService) seedServiceRecords(ctx context.Context, userID uuid.UUID, equipIDs map[string]uuid.UUID) error {
	now := time.Now()
	records := []struct {
		equipName   string
		description string
		cost        float64
		technician  string
		monthsAgo   int
	}{
		{"Sand Filter", "Annual filter clean and inspection", 150, "AquaTech Pool Service", 3},
		{"Variable Speed Pump", "Replaced pump seal and motor bearing", 275, "Pool Pros Inc.", 6},
		{"Gas Heater", "Annual heater tune-up and safety check", 195, "AquaTech Pool Service", 2},
		{"Salt Chlorinator", "Cleaned salt cell, inspected flow sensor", 95, "Self", 1},
		{"Variable Speed Pump", "Routine pump inspection and lubrication", 75, "Pool Pros Inc.", 10},
	}

	for _, r := range records {
		equipID, ok := equipIDs[r.equipName]
		if !ok {
			continue
		}
		serviceDate := now.AddDate(0, -r.monthsAgo, 0)
		record := entities.NewServiceRecord(userID, equipID, serviceDate, r.description, r.cost, r.technician)
		if err := s.srRepo.Create(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (s *DemoSeedService) seedChemicals(ctx context.Context, userID uuid.UUID) error {
	chemicals := []struct {
		name      string
		chemType  entities.ChemicalType
		amount    float64
		unit      valueobjects.Unit
		threshold float64
	}{
		{"Liquid Chlorine", entities.ChemicalTypeSanitizer, 3.5, valueobjects.UnitGallons, 1.0},
		{"pH Decreaser (Muriatic Acid)", entities.ChemicalTypeBalancer, 2.0, valueobjects.UnitGallons, 0.5},
		{"Alkalinity Increaser", entities.ChemicalTypeBalancer, 4.0, valueobjects.UnitPounds, 2.0},
		{"CYA / Stabilizer", entities.ChemicalTypeBalancer, 3.0, valueobjects.UnitPounds, 1.0},
		{"Calcium Hardness Increaser", entities.ChemicalTypeBalancer, 5.0, valueobjects.UnitPounds, 2.0},
	}

	for _, c := range chemicals {
		qty, _ := valueobjects.NewQuantity(c.amount, c.unit)
		chem := entities.NewChemical(userID, c.name, c.chemType, qty, c.threshold)
		if err := s.chemRepo.Create(ctx, chem); err != nil {
			return err
		}
	}
	return nil
}

func randNorm(rng *rand.Rand, mean, stddev float64) float64 {
	return mean + rng.NormFloat64()*stddev
}

func clampF(v, min, max float64) float64 {
	return math.Max(min, math.Min(max, v))
}
