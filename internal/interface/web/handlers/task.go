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

type TaskHandler struct {
	svc *services.TaskService
}

func NewTaskHandler(svc *services.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

type taskSignals struct {
	Name                string `json:"taskName"`
	Description         string `json:"taskDescription"`
	RecurrenceFrequency string `json:"recurrenceFrequency"`
	RecurrenceInterval  int    `json:"recurrenceInterval"`
	DueDate             string `json:"dueDate"`
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(tasks))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *TaskHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	dueDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	sse.PatchElements(renderTaskModal("Add Task", fmt.Sprintf(`
		<div data-signals:taskName="''" data-signals:taskDescription="''"
		     data-signals:recurrenceFrequency="'weekly'" data-signals:recurrenceInterval="1"
		     data-signals:dueDate="'%s'">
			%s
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/tasks')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@post('/tasks')" class="button is-link">Save</button></div>
			</div>
		</div>
	`, dueDate, taskFormFields())))
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	signals := &taskSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dueDate, _ := time.Parse("2006-01-02", signals.DueDate)
	_, err := h.svc.Create(r.Context(), command.CreateTask{
		Name:                signals.Name,
		Description:         signals.Description,
		RecurrenceFrequency: signals.RecurrenceFrequency,
		RecurrenceInterval:  signals.RecurrenceInterval,
		DueDate:             dueDate,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	tasks, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(tasks))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *TaskHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.svc.Get(r.Context(), id)
	if err != nil || task == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderTaskModal("Edit Task", fmt.Sprintf(`
		<div data-signals:taskName="'%s'" data-signals:taskDescription="'%s'"
		     data-signals:recurrenceFrequency="'%s'" data-signals:recurrenceInterval="%d"
		     data-signals:dueDate="'%s'">
			%s
			<div class="field is-grouped is-grouped-right mt-4">
				<div class="control"><button data-on:click="@get('/tasks')" class="button">Cancel</button></div>
				<div class="control"><button data-on:click="@put('/tasks/%s')" class="button is-link">Update</button></div>
			</div>
		</div>
	`, escapeJS(task.Name), escapeJS(task.Description), task.Recurrence.Frequency, task.Recurrence.Interval, task.DueDate.Format("2006-01-02"), taskFormFields(), task.ID.String())))
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	signals := &taskSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dueDate, _ := time.Parse("2006-01-02", signals.DueDate)
	_, err := h.svc.Update(r.Context(), command.UpdateTask{
		ID:                  id,
		Name:                signals.Name,
		Description:         signals.Description,
		RecurrenceFrequency: signals.RecurrenceFrequency,
		RecurrenceInterval:  signals.RecurrenceInterval,
		DueDate:             dueDate,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElements(fmt.Sprintf(`<div id="modal-error" class="notification is-danger is-light">%s</div>`, html.EscapeString(err.Error())))
		return
	}

	tasks, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(tasks))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *TaskHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.svc.Complete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(tasks))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks, _ := h.svc.List(r.Context())
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderList(tasks))
}

func (h *TaskHandler) renderList(tasks []entities.Task) string {
	var b strings.Builder
	b.WriteString(`<div id="tab-content">`)
	b.WriteString(`<div class="level"><div class="level-left"><div class="level-item"><h2 class="title is-4">Maintenance Tasks</h2></div></div>`)
	b.WriteString(`<div class="level-right"><div class="level-item"><button data-on:click="@get('/tasks/new')" class="button is-link">+ Add Task</button></div></div></div>`)

	if len(tasks) == 0 {
		b.WriteString(`<div class="has-text-centered py-6 has-text-grey-light"><p class="is-size-5">No tasks yet</p><p class="is-size-7 mt-1">Add your first maintenance task</p></div>`)
	} else {
		for _, t := range tasks {
			b.WriteString(`<div class="box">`)
			b.WriteString(`<div class="level is-mobile">`)
			b.WriteString(`<div class="level-left">`)
			b.WriteString(fmt.Sprintf(`<div class="level-item">%s</div>`, completeButton(t)))
			b.WriteString(`<div class="level-item"><div>`)
			b.WriteString(fmt.Sprintf(`<p class="has-text-weight-semibold">%s</p>`, html.EscapeString(t.Name)))
			b.WriteString(fmt.Sprintf(`<p class="is-size-7 has-text-grey">Due: %s &middot; Every %d %s</p>`, t.DueDate.Format("Jan 2, 2006"), t.Recurrence.Interval, t.Recurrence.Frequency))
			b.WriteString(`</div></div></div>`)
			b.WriteString(`<div class="level-right"><div class="level-item">`)
			b.WriteString(fmt.Sprintf(`<div class="tags">%s</div>`, taskStatusTag(t.Status)))
			b.WriteString(`</div><div class="level-item"><div class="buttons are-small">`)
			b.WriteString(fmt.Sprintf(`<button data-on:click="@get('/tasks/%s/edit')" class="button is-link is-outlined is-small">Edit</button>`, t.ID.String()))
			b.WriteString(fmt.Sprintf(`<button data-on:click="@delete('/tasks/%s')" class="button is-danger is-outlined is-small">Delete</button>`, t.ID.String()))
			b.WriteString(`</div></div></div>`)
			b.WriteString(`</div></div>`)
		}
	}
	b.WriteString(`</div>`)
	return b.String()
}

func completeButton(t entities.Task) string {
	if t.Status == entities.TaskStatusCompleted {
		return `<span class="icon has-text-success"><i>&#10003;</i></span>`
	}
	return fmt.Sprintf(`<button data-on:click="@post('/tasks/%s/complete')" class="button is-small is-rounded is-white" title="Mark complete" style="width:28px;height:28px;border:2px solid #dbdbdb;padding:0;"></button>`, t.ID.String())
}

func taskStatusTag(status entities.TaskStatus) string {
	switch status {
	case entities.TaskStatusCompleted:
		return `<span class="tag is-success is-light">Completed</span>`
	case entities.TaskStatusOverdue:
		return `<span class="tag is-danger is-light">Overdue</span>`
	default:
		return `<span class="tag is-warning is-light">Pending</span>`
	}
}

func taskFormFields() string {
	return `<div>
		<div class="field"><label class="label">Name</label><div class="control"><input data-bind:taskName type="text" class="input"></div></div>
		<div class="field"><label class="label">Description</label><div class="control"><textarea data-bind:taskDescription rows="2" class="textarea"></textarea></div></div>
		<div class="columns">
			<div class="column">
				<div class="field"><label class="label">Frequency</label><div class="control"><div class="select is-fullwidth"><select data-bind:recurrenceFrequency>
					<option value="daily">Daily</option><option value="weekly">Weekly</option><option value="monthly">Monthly</option>
				</select></div></div></div>
			</div>
			<div class="column">
				<div class="field"><label class="label">Interval</label><div class="control"><input data-bind:recurrenceInterval type="number" min="1" class="input"></div></div>
			</div>
			<div class="column">
				<div class="field"><label class="label">Due Date</label><div class="control"><input data-bind:dueDate type="date" class="input"></div></div>
			</div>
		</div>
	</div>`
}

func renderTaskModal(title, content string) string {
	return fmt.Sprintf(`<div id="modal" class="modal is-active">
		<div class="modal-background" data-on:click="@get('/tasks')"></div>
		<div class="modal-card" style="max-width: 600px;">
			<header class="modal-card-head"><p class="modal-card-title">%s</p></header>
			<section class="modal-card-body">
				<div id="modal-error"></div>
				%s
			</section>
		</div>
	</div>`, title, content)
}
