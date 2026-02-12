package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/josh/poolio/internal/domain/valueobjects"
)

type ChemicalType string

const (
	ChemicalTypeSanitizer ChemicalType = "sanitizer"
	ChemicalTypeShock     ChemicalType = "shock"
	ChemicalTypeBalancer  ChemicalType = "balancer"
	ChemicalTypeAlgaecide ChemicalType = "algaecide"
	ChemicalTypeClarifier ChemicalType = "clarifier"
	ChemicalTypeOther     ChemicalType = "other"
)

type Chemical struct {
	ID             uuid.UUID
	Name           string
	Type           ChemicalType
	Stock          valueobjects.Quantity
	AlertThreshold float64
	LastPurchased  *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewChemical(name string, chemType ChemicalType, stock valueobjects.Quantity, alertThreshold float64) *Chemical {
	now := time.Now()
	return &Chemical{
		ID:             uuid.Must(uuid.NewV7()),
		Name:           name,
		Type:           chemType,
		Stock:          stock,
		AlertThreshold: alertThreshold,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (c *Chemical) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch c.Type {
	case ChemicalTypeSanitizer, ChemicalTypeShock, ChemicalTypeBalancer, ChemicalTypeAlgaecide, ChemicalTypeClarifier, ChemicalTypeOther:
	default:
		return fmt.Errorf("invalid chemical type: %s", c.Type)
	}
	if c.AlertThreshold < 0 {
		return fmt.Errorf("alert threshold cannot be negative")
	}
	return nil
}

func (c *Chemical) IsLowStock() bool {
	return c.Stock.Amount <= c.AlertThreshold
}

func (c *Chemical) AdjustStock(delta float64) error {
	newAmount := c.Stock.Amount + delta
	if newAmount < 0 {
		return fmt.Errorf("stock cannot go below zero")
	}
	c.Stock.Amount = newAmount
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Chemical) RecordPurchase(amount float64) {
	c.Stock.Amount += amount
	now := time.Now()
	c.LastPurchased = &now
	c.UpdatedAt = now
}
