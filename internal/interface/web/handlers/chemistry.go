package handlers

import (
	"net/http"
	"time"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
	"github.com/starfederation/datastar-go/datastar"
)

type ChemistryHandler struct {
	svc     *services.ChemistryService
	userSvc *services.UserService
}

func NewChemistryHandler(svc *services.ChemistryService, userSvc *services.UserService) *ChemistryHandler {
	return &ChemistryHandler{svc: svc, userSvc: userSvc}
}

type chemistrySignals struct {
	PH               float64 `json:"ph"`
	FreeChlorine     float64 `json:"freeChlorine"`
	CombinedChlorine float64 `json:"combinedChlorine"`
	TotalAlkalinity  float64 `json:"totalAlkalinity"`
	CYA              float64 `json:"cya"`
	CalciumHardness  float64 `json:"calciumHardness"`
	Temperature      float64 `json:"temperature"`
	Notes            string  `json:"notes"`
	TestedAt         string  `json:"testedAt"`
}

func (h *ChemistryHandler) List(w http.ResponseWriter, r *http.Request) {
	logs, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemistryList(logs))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *ChemistryHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemistryNewForm(time.Now()))
}

func (h *ChemistryHandler) Create(w http.ResponseWriter, r *http.Request) {
	signals := &chemistrySignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	testedAt, _ := time.Parse("2006-01-02T15:04", signals.TestedAt)
	_, err := h.svc.Create(r.Context(), command.CreateChemistryLog{
		PH:               signals.PH,
		FreeChlorine:     signals.FreeChlorine,
		CombinedChlorine: signals.CombinedChlorine,
		TotalAlkalinity:  signals.TotalAlkalinity,
		CYA:              signals.CYA,
		CalciumHardness:  signals.CalciumHardness,
		Temperature:      signals.Temperature,
		Notes:            signals.Notes,
		TestedAt:         testedAt,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError(err.Error()))
		return
	}

	logs, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemistryList(logs))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *ChemistryHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log, err := h.svc.Get(r.Context(), id)
	if err != nil || log == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemistryEditForm(log))
}

func (h *ChemistryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	signals := &chemistrySignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	testedAt, _ := time.Parse("2006-01-02T15:04", signals.TestedAt)
	_, err := h.svc.Update(r.Context(), command.UpdateChemistryLog{
		ID:               id,
		PH:               signals.PH,
		FreeChlorine:     signals.FreeChlorine,
		CombinedChlorine: signals.CombinedChlorine,
		TotalAlkalinity:  signals.TotalAlkalinity,
		CYA:              signals.CYA,
		CalciumHardness:  signals.CalciumHardness,
		Temperature:      signals.Temperature,
		Notes:            signals.Notes,
		TestedAt:         testedAt,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError(err.Error()))
		return
	}

	logs, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemistryList(logs))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *ChemistryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ChemistryList(logs))
}

func (h *ChemistryHandler) Plan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log, err := h.svc.Get(r.Context(), id)
	if err != nil || log == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	user, err := services.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	plan := entities.GenerateTreatmentPlan(log, user.PoolGallons)

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TreatmentPlanModal(plan))
}
