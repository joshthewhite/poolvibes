package handlers

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/josh/poolio/internal/application/command"
	"github.com/josh/poolio/internal/application/services"
	"github.com/josh/poolio/internal/domain/entities"
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
	sse.PatchElements(h.renderList(items))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *EquipmentHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderEquipModal("Add Equipment", `
		<div data-signals:eqName="''" data-signals:eqCategory="'pump'" data-signals:eqManufacturer="''"
		     data-signals:eqModel="''" data-signals:eqSerialNumber="''" data-signals:eqInstallDate="''" data-signals:eqWarrantyExpiry="''">
			`+equipmentFormFields()+`
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/equipment')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@post('/equipment')" class="button is-link">Save</button></div>
			</div>
		</div>
	`))
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
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(items))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *EquipmentHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	eq, err := h.svc.Get(r.Context(), id)
	if err != nil || eq == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderEquipModal("Edit Equipment", fmt.Sprintf(`
		<div data-signals:eqName="'%s'" data-signals:eqCategory="'%s'" data-signals:eqManufacturer="'%s'"
		     data-signals:eqModel="'%s'" data-signals:eqSerialNumber="'%s'" data-signals:eqInstallDate="'%s'" data-signals:eqWarrantyExpiry="'%s'">
			%s
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/equipment')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@put('/equipment/%s')" class="button is-link">Update</button></div>
			</div>
		</div>
	`, escapeJS(eq.Name), eq.Category, escapeJS(eq.Manufacturer), escapeJS(eq.Model), escapeJS(eq.SerialNumber), fmtDatePtr(eq.InstallDate), fmtDatePtr(eq.WarrantyExpiry), equipmentFormFields(), eq.ID.String())))
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
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(items))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *EquipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(items))
}

func (h *EquipmentHandler) NewServiceRecordForm(w http.ResponseWriter, r *http.Request) {
	eqID := r.PathValue("id")
	today := time.Now().Format("2006-01-02")

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderEquipModal("Add Service Record", fmt.Sprintf(`
		<div data-signals:srServiceDate="'%s'" data-signals:srDescription="''" data-signals:srCost="0" data-signals:srTechnician="''">
			<div class="field"><label class="label">Service Date</label><div class="control"><input data-bind:srServiceDate type="date" class="input"></div></div>
			<div class="field"><label class="label">Description</label><div class="control"><textarea data-bind:srDescription rows="2" class="textarea"></textarea></div></div>
			<div class="columns">
				<div class="column"><div class="field"><label class="label">Cost ($)</label><div class="control"><input data-bind:srCost type="number" step="0.01" min="0" class="input"></div></div></div>
				<div class="column"><div class="field"><label class="label">Technician</label><div class="control"><input data-bind:srTechnician type="text" class="input"></div></div></div>
			</div>
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/equipment')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@post('/equipment/%s/service-records')" class="button is-link">Save</button></div>
			</div>
		</div>
	`, today, eqID)))
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
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(items))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *EquipmentHandler) DeleteServiceRecord(w http.ResponseWriter, r *http.Request) {
	recordID := r.PathValue("recordId")
	if err := h.svc.DeleteServiceRecord(r.Context(), recordID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(items))
}

func (h *EquipmentHandler) renderList(items []entities.Equipment) string {
	var b strings.Builder
	b.WriteString(`<div id="tab-content">`)
	b.WriteString(`<div class="level"><div class="level-left"><div class="level-item"><h2 class="title is-4">Equipment</h2></div></div>`)
	b.WriteString(`<div class="level-right"><div class="level-item"><button data-on:click="@get('/equipment/new')" class="button is-link">+ Add Equipment</button></div></div></div>`)

	if len(items) == 0 {
		b.WriteString(`<div class="has-text-centered py-6 has-text-grey-light"><p class="is-size-5">No equipment yet</p><p class="is-size-7 mt-1">Add your pool equipment to track maintenance</p></div>`)
	} else {
		b.WriteString(`<div class="columns is-multiline">`)
		for _, eq := range items {
			b.WriteString(`<div class="column is-half">`)
			b.WriteString(`<div class="card">`)
			b.WriteString(`<div class="card-content">`)

			// Header
			b.WriteString(`<div class="level is-mobile mb-3">`)
			b.WriteString(fmt.Sprintf(`<div class="level-left"><div class="level-item"><div><p class="title is-5 mb-1">%s</p><span class="tag is-light">%s</span></div></div></div>`, html.EscapeString(eq.Name), eq.Category))
			b.WriteString(fmt.Sprintf(`<div class="level-right"><div class="level-item"><div class="buttons are-small"><button data-on:click="@get('/equipment/%s/edit')" class="button is-link is-outlined is-small">Edit</button><button data-on:click="@delete('/equipment/%s')" class="button is-danger is-outlined is-small">Delete</button></div></div></div>`, eq.ID.String(), eq.ID.String()))
			b.WriteString(`</div>`)

			if eq.Manufacturer != "" || eq.Model != "" {
				b.WriteString(fmt.Sprintf(`<p class="is-size-7 has-text-grey">%s %s</p>`, html.EscapeString(eq.Manufacturer), html.EscapeString(eq.Model)))
			}
			if eq.SerialNumber != "" {
				b.WriteString(fmt.Sprintf(`<p class="is-size-7 has-text-grey-light">S/N: %s</p>`, html.EscapeString(eq.SerialNumber)))
			}
			if eq.WarrantyExpiry != nil {
				warrantyTag := "is-success"
				if !eq.IsWarrantyActive() {
					warrantyTag = "is-danger"
				}
				b.WriteString(fmt.Sprintf(`<p class="is-size-7 mt-1"><span class="tag %s is-light is-small">Warranty: %s</span></p>`, warrantyTag, eq.WarrantyExpiry.Format("Jan 2, 2006")))
			}

			// Service records
			b.WriteString(`<hr class="my-3" style="height:1px;background:#f0f0f0;border:none;">`)
			b.WriteString(fmt.Sprintf(`<div class="level is-mobile mb-2"><div class="level-left"><div class="level-item"><span class="has-text-weight-semibold is-size-7">Service History</span></div></div><div class="level-right"><div class="level-item"><button data-on:click="@get('/equipment/%s/service-records/new')" class="button is-link is-outlined is-small">+ Add</button></div></div></div>`, eq.ID.String()))

			if len(eq.ServiceRecords) == 0 {
				b.WriteString(`<p class="is-size-7 has-text-grey-light">No service records</p>`)
			} else {
				for _, sr := range eq.ServiceRecords {
					b.WriteString(`<div class="level is-mobile mb-1" style="margin-bottom:0.25rem!important;">`)
					b.WriteString(fmt.Sprintf(`<div class="level-left"><div class="level-item"><span class="is-size-7">%s</span></div><div class="level-item"><span class="is-size-7 has-text-grey-light">%s</span></div></div>`, html.EscapeString(sr.Description), sr.ServiceDate.Format("Jan 2")))
					b.WriteString(`<div class="level-right">`)
					if sr.Cost > 0 {
						b.WriteString(fmt.Sprintf(`<div class="level-item"><span class="is-size-7 has-text-grey">$%.2f</span></div>`, sr.Cost))
					}
					b.WriteString(fmt.Sprintf(`<div class="level-item"><button data-on:click="@delete('/equipment/%s/service-records/%s')" class="delete is-small"></button></div>`, eq.ID.String(), sr.ID.String()))
					b.WriteString(`</div></div>`)
				}
			}

			b.WriteString(`</div></div></div>`)
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func equipmentFormFields() string {
	return `<div>
		<div class="field"><label class="label">Name</label><div class="control"><input data-bind:eqName type="text" class="input"></div></div>
		<div class="field"><label class="label">Category</label><div class="control"><div class="select is-fullwidth"><select data-bind:eqCategory>
			<option value="pump">Pump</option><option value="filter">Filter</option><option value="heater">Heater</option><option value="chlorinator">Chlorinator</option><option value="cleaner">Cleaner</option><option value="other">Other</option>
		</select></div></div></div>
		<div class="columns">
			<div class="column"><div class="field"><label class="label">Manufacturer</label><div class="control"><input data-bind:eqManufacturer type="text" class="input"></div></div></div>
			<div class="column"><div class="field"><label class="label">Model</label><div class="control"><input data-bind:eqModel type="text" class="input"></div></div></div>
		</div>
		<div class="field"><label class="label">Serial Number</label><div class="control"><input data-bind:eqSerialNumber type="text" class="input"></div></div>
		<div class="columns">
			<div class="column"><div class="field"><label class="label">Install Date</label><div class="control"><input data-bind:eqInstallDate type="date" class="input"></div></div></div>
			<div class="column"><div class="field"><label class="label">Warranty Expiry</label><div class="control"><input data-bind:eqWarrantyExpiry type="date" class="input"></div></div></div>
		</div>
	</div>`
}

func renderEquipModal(title, content string) string {
	return fmt.Sprintf(`<div id="modal" class="modal is-active">
		<div class="modal-background" data-on:click="@get('/equipment')"></div>
		<div class="modal-card" style="max-width: 600px;">
			<header class="modal-card-head"><p class="modal-card-title">%s</p></header>
			<section class="modal-card-body">
				<div id="modal-error"></div>
				%s
			</section>
		</div>
	</div>`, title, content)
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

func fmtDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}
