package command

import "time"

type CreateEquipment struct {
	Name           string
	Category       string
	Manufacturer   string
	Model          string
	SerialNumber   string
	InstallDate    *time.Time
	WarrantyExpiry *time.Time
}

type UpdateEquipment struct {
	ID             string
	Name           string
	Category       string
	Manufacturer   string
	Model          string
	SerialNumber   string
	InstallDate    *time.Time
	WarrantyExpiry *time.Time
}

type CreateServiceRecord struct {
	EquipmentID string
	ServiceDate time.Time
	Description string
	Cost        float64
	Technician  string
}
