package handlers

import (
	"log/slog"
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
	"github.com/starfederation/datastar-go/datastar"
)

type AdminHandler struct {
	svc *services.UserService
}

func NewAdminHandler(svc *services.UserService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

type adminUserSignals struct {
	IsAdmin    bool `json:"isAdmin"`
	IsDisabled bool `json:"isDisabled"`
	IsDemo     bool `json:"isDemo"`
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		slog.Error("Error listing users", "error", err)
		http.Error(w, "failed to load users", http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.AdminUserList(users))
	sse.PatchElementTempl(templates.EmptyModal())
}

func (h *AdminHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.svc.Get(r.Context(), id)
	if err != nil {
		slog.Error("Error loading user", "userID", id, "error", err)
		http.Error(w, "failed to load user", http.StatusBadRequest)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.AdminEditUser(user))
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var signals adminUserSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		http.Error(w, "invalid request data", http.StatusBadRequest)
		return
	}

	_, err := h.svc.Update(r.Context(), command.UpdateUser{
		ID:         id,
		IsAdmin:    signals.IsAdmin,
		IsDisabled: signals.IsDisabled,
		IsDemo:     signals.IsDemo,
	})
	if err != nil {
		msg := err.Error()
		if msg != "cannot modify your own account" && msg != "cannot remove the last admin" {
			slog.Error("Error updating user", "userID", id, "error", err)
			msg = "failed to update user"
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	users, err := h.svc.List(r.Context())
	if err != nil {
		slog.Error("Error listing users", "error", err)
		http.Error(w, "failed to load users", http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.AdminUserList(users))
	sse.PatchElementTempl(templates.EmptyModal())
}
