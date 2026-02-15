package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
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
		slog.Error("Error listing tasks", "error", err)
		http.Error(w, "failed to load tasks", http.StatusInternalServerError)
		return
	}

	active, completed := splitTasks(tasks)
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TaskList(active, completed))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *TaskHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	dueDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TaskNewForm(dueDate))
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	signals := &taskSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, "invalid request data", http.StatusBadRequest)
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
		slog.Error("Error creating task", "error", err)
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError("Failed to create task"))
		return
	}

	tasks, _ := h.svc.List(r.Context())
	active, completed := splitTasks(tasks)
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TaskList(active, completed))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *TaskHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.svc.Get(r.Context(), id)
	if err != nil || task == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TaskEditForm(task))
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	signals := &taskSignals{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, "invalid request data", http.StatusBadRequest)
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
		slog.Error("Error updating task", "error", err)
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.ModalError("Failed to update task"))
		return
	}

	tasks, _ := h.svc.List(r.Context())
	active, completed := splitTasks(tasks)
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TaskList(active, completed))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *TaskHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := h.svc.Complete(r.Context(), id)
	if err != nil {
		slog.Error("Error completing task", "error", err)
		http.Error(w, "failed to complete task", http.StatusInternalServerError)
		return
	}

	tasks, _ := h.svc.List(r.Context())
	active, completed := splitTasks(tasks)
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TaskList(active, completed))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		slog.Error("Error deleting task", "error", err)
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}

	tasks, _ := h.svc.List(r.Context())
	active, completed := splitTasks(tasks)
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.TaskList(active, completed))
}

func splitTasks(tasks []entities.Task) (active, completed []entities.Task) {
	for _, t := range tasks {
		if t.Status == entities.TaskStatusCompleted {
			completed = append(completed, t)
		} else {
			active = append(active, t)
		}
	}
	return
}
