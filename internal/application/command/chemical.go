package command

type CreateChemical struct {
	Name           string
	Type           string
	StockAmount    float64
	StockUnit      string
	AlertThreshold float64
}

type UpdateChemical struct {
	ID             string
	Name           string
	Type           string
	StockAmount    float64
	StockUnit      string
	AlertThreshold float64
}

type AdjustChemicalStock struct {
	ID    string
	Delta float64
}
