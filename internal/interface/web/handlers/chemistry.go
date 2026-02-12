package handlers

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/josh/poolio/internal/application/command"
	"github.com/josh/poolio/internal/application/services"
	"github.com/josh/poolio/internal/domain/entities"
	"github.com/starfederation/datastar-go/datastar"
)

type ChemistryHandler struct {
	svc *services.ChemistryService
}

func NewChemistryHandler(svc *services.ChemistryService) *ChemistryHandler {
	return &ChemistryHandler{svc: svc}
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
	sse.PatchElements(h.renderList(logs))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *ChemistryHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	now := time.Now().Format("2006-01-02T15:04")
	sse.PatchElements(renderModal("Add Chemistry Log", fmt.Sprintf(`
		<div data-signals:ph="7.4" data-signals:freeChlorine="2.0" data-signals:combinedChlorine="0.0"
		     data-signals:totalAlkalinity="100" data-signals:cya="40" data-signals:calciumHardness="300"
		     data-signals:temperature="80" data-signals:notes="''" data-signals:testedAt="'%s'">
			<div class="columns is-multiline">
				<div class="column is-half">
					<div class="field"><label class="label">pH</label><div class="control"><input data-bind:ph type="number" step="0.1" min="0" max="14" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Free Chlorine (ppm)</label><div class="control"><input data-bind:freeChlorine type="number" step="0.1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Combined Chlorine (ppm)</label><div class="control"><input data-bind:combinedChlorine type="number" step="0.1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Total Alkalinity (ppm)</label><div class="control"><input data-bind:totalAlkalinity type="number" step="1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">CYA (ppm)</label><div class="control"><input data-bind:cya type="number" step="1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Calcium Hardness (ppm)</label><div class="control"><input data-bind:calciumHardness type="number" step="1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Temperature (&deg;F)</label><div class="control"><input data-bind:temperature type="number" step="1" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Tested At</label><div class="control"><input data-bind:testedAt type="datetime-local" class="input"></div></div>
				</div>
				<div class="column is-full">
					<div class="field"><label class="label">Notes</label><div class="control"><textarea data-bind:notes rows="2" class="textarea"></textarea></div></div>
				</div>
			</div>
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/chemistry')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@post('/chemistry')" class="button is-link">Save</button></div>
			</div>
		</div>
	`, now)))
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
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	logs, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(logs))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *ChemistryHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log, err := h.svc.Get(r.Context(), id)
	if err != nil || log == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderModal("Edit Chemistry Log", fmt.Sprintf(`
		<div data-signals:ph="%g" data-signals:freeChlorine="%g" data-signals:combinedChlorine="%g"
		     data-signals:totalAlkalinity="%g" data-signals:cya="%g" data-signals:calciumHardness="%g"
		     data-signals:temperature="%g" data-signals:notes="'%s'" data-signals:testedAt="'%s'">
			<div class="columns is-multiline">
				<div class="column is-half">
					<div class="field"><label class="label">pH</label><div class="control"><input data-bind:ph type="number" step="0.1" min="0" max="14" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Free Chlorine (ppm)</label><div class="control"><input data-bind:freeChlorine type="number" step="0.1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Combined Chlorine (ppm)</label><div class="control"><input data-bind:combinedChlorine type="number" step="0.1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Total Alkalinity (ppm)</label><div class="control"><input data-bind:totalAlkalinity type="number" step="1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">CYA (ppm)</label><div class="control"><input data-bind:cya type="number" step="1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Calcium Hardness (ppm)</label><div class="control"><input data-bind:calciumHardness type="number" step="1" min="0" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Temperature (&deg;F)</label><div class="control"><input data-bind:temperature type="number" step="1" class="input"></div></div>
				</div>
				<div class="column is-half">
					<div class="field"><label class="label">Tested At</label><div class="control"><input data-bind:testedAt type="datetime-local" class="input"></div></div>
				</div>
				<div class="column is-full">
					<div class="field"><label class="label">Notes</label><div class="control"><textarea data-bind:notes rows="2" class="textarea"></textarea></div></div>
				</div>
			</div>
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/chemistry')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@put('/chemistry/%s')" class="button is-link">Update</button></div>
			</div>
		</div>
	`, log.PH, log.FreeChlorine, log.CombinedChlorine, log.TotalAlkalinity, log.CYA, log.CalciumHardness, log.Temperature, escapeJS(log.Notes), log.TestedAt.Format("2006-01-02T15:04"), log.ID.String())))
}

func (h *ChemistryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	logs, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(logs))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *ChemistryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(logs))
}

func (h *ChemistryHandler) renderList(logs []entities.ChemistryLog) string {
	var b strings.Builder
	b.WriteString(`<div id="tab-content">`)
	b.WriteString(`<div class="level"><div class="level-left"><div class="level-item"><h2 class="title is-4">Water Chemistry Logs</h2></div></div>`)
	b.WriteString(`<div class="level-right"><div class="level-item"><button data-on:click="@get('/chemistry/new')" class="button is-link">+ Add Test</button></div></div></div>`)

	if len(logs) == 0 {
		b.WriteString(`<div class="has-text-centered py-6 has-text-grey-light"><p class="is-size-5">No chemistry logs yet</p><p class="is-size-7 mt-1">Add your first water test to get started</p></div>`)
	} else {
		b.WriteString(`<div class="table-container"><table class="table is-fullwidth is-hoverable is-striped">`)
		b.WriteString(`<thead><tr><th>Date</th><th>pH</th><th>FC</th><th>CC</th><th>TA</th><th>CYA</th><th>CH</th><th>Temp</th><th class="has-text-right">Actions</th></tr></thead><tbody>`)

		for _, l := range logs {
			b.WriteString(`<tr>`)
			b.WriteString(fmt.Sprintf(`<td>%s</td>`, l.TestedAt.Format("Jan 2, 3:04 PM")))
			b.WriteString(fmt.Sprintf(`<td><span class="%s">%.1f</span></td>`, valueClass(l.PHInRange()), l.PH))
			b.WriteString(fmt.Sprintf(`<td><span class="%s">%.1f</span></td>`, valueClass(l.FreeChlorineInRange()), l.FreeChlorine))
			b.WriteString(fmt.Sprintf(`<td><span class="%s">%.1f</span></td>`, valueClass(l.CombinedChlorineInRange()), l.CombinedChlorine))
			b.WriteString(fmt.Sprintf(`<td><span class="%s">%.0f</span></td>`, valueClass(l.TotalAlkalinityInRange()), l.TotalAlkalinity))
			b.WriteString(fmt.Sprintf(`<td><span class="%s">%.0f</span></td>`, valueClass(l.CYAInRange()), l.CYA))
			b.WriteString(fmt.Sprintf(`<td><span class="%s">%.0f</span></td>`, valueClass(l.CalciumHardnessInRange()), l.CalciumHardness))
			b.WriteString(fmt.Sprintf(`<td>%.0f&deg;F</td>`, l.Temperature))
			b.WriteString(fmt.Sprintf(`<td class="has-text-right"><div class="buttons is-right are-small"><button data-on:click="@get('/chemistry/%s/edit')" class="button is-link is-outlined is-small">Edit</button><button data-on:click="@delete('/chemistry/%s')" class="button is-danger is-outlined is-small">Delete</button></div></td>`, l.ID.String(), l.ID.String()))
			b.WriteString(`</tr>`)
		}
		b.WriteString(`</tbody></table></div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func valueClass(inRange bool) string {
	if inRange {
		return "has-text-success has-text-weight-semibold"
	}
	return "has-text-danger has-text-weight-bold"
}

func escapeJS(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}

func renderModal(title, content string) string {
	return fmt.Sprintf(`<div id="modal" class="modal is-active">
		<div class="modal-background" data-on:click="@get('/chemistry')"></div>
		<div class="modal-card" style="max-width: 720px;">
			<header class="modal-card-head"><p class="modal-card-title">%s</p></header>
			<section class="modal-card-body">
				<div id="modal-error"></div>
				%s
			</section>
		</div>
	</div>`, title, content)
}
