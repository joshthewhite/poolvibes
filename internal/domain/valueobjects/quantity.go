package valueobjects

import "fmt"

type Unit string

const (
	UnitPounds  Unit = "lbs"
	UnitOunces  Unit = "oz"
	UnitGallons Unit = "gal"
	UnitLiters  Unit = "L"
	UnitKg      Unit = "kg"
)

type Quantity struct {
	Amount float64
	Unit   Unit
}

func NewQuantity(amount float64, unit Unit) (Quantity, error) {
	if amount < 0 {
		return Quantity{}, fmt.Errorf("amount cannot be negative")
	}
	switch unit {
	case UnitPounds, UnitOunces, UnitGallons, UnitLiters, UnitKg:
	default:
		return Quantity{}, fmt.Errorf("invalid unit: %s", unit)
	}
	return Quantity{Amount: amount, Unit: unit}, nil
}
