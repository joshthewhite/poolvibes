package handlers

import (
	"net/http"
	"time"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
	"github.com/starfederation/datastar-go/datastar"
)

type EquipmentHandler struct {
	svc *services.EquipmentService
}

func NewEquipmentHandler(svc *services.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{svc: svc}
}

type equipmentSignals struct {
	Name           string `json:"eqName"`
	Category       string `json:"eqCategory"`
	Manufacturer   string `json:"eqManufacturer"`
	Model          string `json:"eqModel"`
	SerialNumber   string `json:"eqSerialNumber"`
	InstallDate    string `json:"eqInstallDate"`
	WarrantyExpiry string `json:"eqWarrantyExpiry"`
}

type serviceRecordSignals struct {
	ServiceDate string  `json:"srServiceDate"`
	Description string  `json:"srDescription"`
	Cost        float64 `json:"srCost"`
	Technician  string  `json:"srTechnician"`
}

func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentList(items))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *EquipmentHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentNewForm())
}

func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	signals := &equipmentSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.CreateEquipment{
		Name:           signals.Name,
		Category:       signals.Category,
		Manufacturer:   signals.Manufacturer,
		Model:          signals.Model,
		SerialNumber:   signals.SerialNumber,
		InstallDate:    parseOptionalDate(signals.InstallDate),
		WarrantyExpiry: parseOptionalDate(signals.WarrantyExpiry),
	}

	_, err := h.svc.Create(r.Context(), cmd)
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError(err.Error()))
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentList(items))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *EquipmentHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	eq, err := h.svc.Get(r.Context(), id)
	if err != nil || eq == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentEditForm(eq))
}

func (h *EquipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	signals := &equipmentSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := command.UpdateEquipment{
		ID:             id,
		Name:           signals.Name,
		Category:       signals.Category,
		Manufacturer:   signals.Manufacturer,
		Model:          signals.Model,
		SerialNumber:   signals.SerialNumber,
		InstallDate:    parseOptionalDate(signals.InstallDate),
		WarrantyExpiry: parseOptionalDate(signals.WarrantyExpiry),
	}

	_, err := h.svc.Update(r.Context(), cmd)
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError(err.Error()))
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentList(items))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *EquipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentList(items))
}

func (h *EquipmentHandler) NewServiceRecordForm(w http.ResponseWriter, r *http.Request) {
	eqID := r.PathValue("id")
	today := time.Now().Format("2006-01-02")

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.ServiceRecordNewForm(eqID, today))
}

func (h *EquipmentHandler) CreateServiceRecord(w http.ResponseWriter, r *http.Request) {
	eqID := r.PathValue("id")
	signals := &serviceRecordSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serviceDate, _ := time.Parse("2006-01-02", signals.ServiceDate)
	_, err := h.svc.AddServiceRecord(r.Context(), command.CreateServiceRecord{
		EquipmentID: eqID,
		ServiceDate: serviceDate,
		Description: signals.Description,
		Cost:        signals.Cost,
		Technician:  signals.Technician,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError(err.Error()))
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentList(items))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *EquipmentHandler) DeleteServiceRecord(w http.ResponseWriter, r *http.Request) {
	recordID := r.PathValue("recordId")
	if err := h.svc.DeleteServiceRecord(r.Context(), recordID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.EquipmentList(items))
}

func parseOptionalDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}
