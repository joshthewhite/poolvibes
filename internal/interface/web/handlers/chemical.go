package handlers

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/josh/poolio/internal/application/command"
	"github.com/josh/poolio/internal/application/services"
	"github.com/josh/poolio/internal/domain/entities"
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
	sse.PatchElements(h.renderList(chemicals))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *ChemicalHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderChemModal("Add Chemical", `
		<div data-signals:chemName="''" data-signals:chemType="'sanitizer'" data-signals:chemStockAmount="0"
		     data-signals:chemStockUnit="'lbs'" data-signals:chemAlertThreshold="5">
			`+chemicalFormFields()+`
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/chemicals')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@post('/chemicals')" class="button is-link">Save</button></div>
			</div>
		</div>
	`))
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
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	chemicals, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(chemicals))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *ChemicalHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	chem, err := h.svc.Get(r.Context(), id)
	if err != nil || chem == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderChemModal("Edit Chemical", fmt.Sprintf(`
		<div data-signals:chemName="'%s'" data-signals:chemType="'%s'" data-signals:chemStockAmount="%g"
		     data-signals:chemStockUnit="'%s'" data-signals:chemAlertThreshold="%g">
			%s
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/chemicals')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@put('/chemicals/%s')" class="button is-link">Update</button></div>
			</div>
		</div>
	`, escapeJS(chem.Name), chem.Type, chem.Stock.Amount, chem.Stock.Unit, chem.AlertThreshold, chemicalFormFields(), chem.ID.String())))
}

func (h *ChemicalHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	chemicals, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(chemicals))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *ChemicalHandler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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
	sse.PatchElements(h.renderList(chemicals))
}

func (h *ChemicalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chemicals, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(chemicals))
}

func (h *ChemicalHandler) renderList(chemicals []entities.Chemical) string {
	var b strings.Builder
	b.WriteString(`<div id="tab-content">`)
	b.WriteString(`<div class="level"><div class="level-left"><div class="level-item"><h2 class="title is-4">Chemical Inventory</h2></div></div>`)
	b.WriteString(`<div class="level-right"><div class="level-item"><button data-on:click="@get('/chemicals/new')" class="button is-link">+ Add Chemical</button></div></div></div>`)

	if len(chemicals) == 0 {
		b.WriteString(`<div class="has-text-centered py-6 has-text-grey-light"><p class="is-size-5">No chemicals tracked</p><p class="is-size-7 mt-1">Add chemicals to track your inventory</p></div>`)
	} else {
		b.WriteString(`<div class="columns is-multiline">`)
		for _, c := range chemicals {
			lowStockNotif := ""
			cardClass := "card"
			if c.IsLowStock() {
				lowStockNotif = `<span class="tag is-danger is-light ml-2">Low Stock</span>`
				cardClass = "card" // border handled below
			}

			b.WriteString(fmt.Sprintf(`<div class="column is-one-third"><div class="%s"`, cardClass))
			if c.IsLowStock() {
				b.WriteString(` style="border: 2px solid hsl(348, 86%%, 61%%); background: hsl(348, 86%%, 97%%);"`)
			}
			b.WriteString(`>`)
			b.WriteString(`<div class="card-content">`)

			// Header
			b.WriteString(fmt.Sprintf(`<div class="level is-mobile mb-2"><div class="level-left"><div class="level-item"><div><p class="has-text-weight-bold">%s</p><div><span class="tag is-light is-small">%s</span>%s</div></div></div></div>`, html.EscapeString(c.Name), c.Type, lowStockNotif))
			b.WriteString(fmt.Sprintf(`<div class="level-right"><div class="level-item"><div class="buttons are-small"><button data-on:click="@get('/chemicals/%s/edit')" class="button is-link is-outlined is-small">Edit</button><button data-on:click="@delete('/chemicals/%s')" class="button is-danger is-outlined is-small">Delete</button></div></div></div></div>`, c.ID.String(), c.ID.String()))

			// Stock display
			b.WriteString(fmt.Sprintf(`<p class="is-size-3 has-text-weight-bold mt-2">%.1f <span class="is-size-6 has-text-weight-normal has-text-grey">%s</span></p>`, c.Stock.Amount, c.Stock.Unit))
			b.WriteString(fmt.Sprintf(`<p class="is-size-7 has-text-grey-light">Alert threshold: %.1f %s</p>`, c.AlertThreshold, c.Stock.Unit))

			// Quick adjust buttons
			b.WriteString(`<hr class="my-3" style="height:1px;background:#f0f0f0;border:none;">`)
			b.WriteString(fmt.Sprintf(`<div class="buttons are-small" data-signals:adjustDelta="0">
				<button data-on:click="$adjustDelta = -1; @post('/chemicals/%s/adjust')" class="button is-small">-1</button>
				<button data-on:click="$adjustDelta = -5; @post('/chemicals/%s/adjust')" class="button is-small">-5</button>
				<button data-on:click="$adjustDelta = 5; @post('/chemicals/%s/adjust')" class="button is-small is-success is-outlined">+5</button>
				<button data-on:click="$adjustDelta = 10; @post('/chemicals/%s/adjust')" class="button is-small is-success is-outlined">+10</button>
			</div>`, c.ID.String(), c.ID.String(), c.ID.String(), c.ID.String()))

			if c.LastPurchased != nil {
				b.WriteString(fmt.Sprintf(`<p class="is-size-7 has-text-grey-light mt-2">Last purchased: %s</p>`, c.LastPurchased.Format("Jan 2, 2006")))
			}

			b.WriteString(`</div></div></div>`)
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func chemicalFormFields() string {
	return `<div>
		<div class="field"><label class="label">Name</label><div class="control"><input data-bind:chemName type="text" class="input"></div></div>
		<div class="field"><label class="label">Type</label><div class="control"><div class="select is-fullwidth"><select data-bind:chemType>
			<option value="sanitizer">Sanitizer</option><option value="shock">Shock</option><option value="balancer">Balancer</option><option value="algaecide">Algaecide</option><option value="clarifier">Clarifier</option><option value="other">Other</option>
		</select></div></div></div>
		<div class="columns">
			<div class="column"><div class="field"><label class="label">Stock Amount</label><div class="control"><input data-bind:chemStockAmount type="number" step="0.1" min="0" class="input"></div></div></div>
			<div class="column"><div class="field"><label class="label">Unit</label><div class="control"><div class="select is-fullwidth"><select data-bind:chemStockUnit>
				<option value="lbs">Pounds (lbs)</option><option value="oz">Ounces (oz)</option><option value="gal">Gallons (gal)</option><option value="L">Liters (L)</option><option value="kg">Kilograms (kg)</option>
			</select></div></div></div></div>
			<div class="column"><div class="field"><label class="label">Alert At</label><div class="control"><input data-bind:chemAlertThreshold type="number" step="0.1" min="0" class="input"></div></div></div>
		</div>
	</div>`
}

func renderChemModal(title, content string) string {
	return fmt.Sprintf(`<div id="modal" class="modal is-active">
		<div class="modal-background" data-on:click="@get('/chemicals')"></div>
		<div class="modal-card" style="max-width: 600px;">
			<header class="modal-card-head"><p class="modal-card-title">%s</p></header>
			<section class="modal-card-body">
				<div id="modal-error"></div>
				%s
			</section>
		</div>
	</div>`, title, content)
}
