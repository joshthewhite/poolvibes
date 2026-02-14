package handlers

import (
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
	"github.com/starfederation/datastar-go/datastar"
)

type ChemicalHandler struct {
	svc *services.ChemicalService
}

func NewChemicalHandler(svc *services.ChemicalService) *ChemicalHandler {
	return &ChemicalHandler{svc: svc}
}

type chemicalSignals struct {
	Name           string  `json:"chemName"`
	Type           string  `json:"chemType"`
	StockAmount    float64 `json:"chemStockAmount"`
	StockUnit      string  `json:"chemStockUnit"`
	AlertThreshold float64 `json:"chemAlertThreshold"`
}

type adjustSignals struct {
	Delta float64 `json:"adjustDelta"`
}

func (h *ChemicalHandler) List(w http.ResponseWriter, r *http.Request) {
	chemicals, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemicalList(chemicals))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *ChemicalHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemicalNewForm())
}

func (h *ChemicalHandler) Create(w http.ResponseWriter, r *http.Request) {
	signals := &chemicalSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.svc.Create(r.Context(), command.CreateChemical{
		Name:           signals.Name,
		Type:           signals.Type,
		StockAmount:    signals.StockAmount,
		StockUnit:      signals.StockUnit,
		AlertThreshold: signals.AlertThreshold,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError(err.Error()))
		return
	}

	chemicals, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemicalList(chemicals))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *ChemicalHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	chem, err := h.svc.Get(r.Context(), id)
	if err != nil || chem == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemicalEditForm(chem))
}

func (h *ChemicalHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	signals := &chemicalSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.svc.Update(r.Context(), command.UpdateChemical{
		ID:             id,
		Name:           signals.Name,
		Type:           signals.Type,
		StockAmount:    signals.StockAmount,
		StockUnit:      signals.StockUnit,
		AlertThreshold: signals.AlertThreshold,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError(err.Error()))
		return
	}

	chemicals, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemicalList(chemicals))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *ChemicalHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	signals := &adjustSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.svc.AdjustStock(r.Context(), command.AdjustChemicalStock{
		ID:    id,
		Delta: signals.Delta,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chemicals, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemicalList(chemicals))
}

func (h *ChemicalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chemicals, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemicalList(chemicals))
}
