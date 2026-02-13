package handlers

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/domain/entities"
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
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderUserList(users))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *AdminHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	isAdmin := "false"
	if user.IsAdmin {
		isAdmin = "true"
	}
	isDisabled := "false"
	if user.IsDisabled {
		isDisabled = "true"
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(renderModal("Edit User", fmt.Sprintf(`
		<div data-signals:isAdmin="%s" data-signals:isDisabled="%s">
			<div class="field">
				<label class="label">Email</label>
				<div class="control">
					<input type="text" class="input" value="%s" disabled>
				</div>
			</div>
			<div class="field">
				<label class="checkbox">
					<input data-bind:isAdmin type="checkbox"> Admin
				</label>
			</div>
			<div class="field">
				<label class="checkbox">
					<input data-bind:isDisabled type="checkbox"> Disabled
				</label>
			</div>
			<div class="field">
				<div class="control">
					<button class="button is-link" data-on:click="@put('/admin/users/%s')">Save</button>
				</div>
			</div>
		</div>`,
		isAdmin,
		isDisabled,
		html.EscapeString(user.Email),
		html.EscapeString(user.ID.String()),
	)))
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var signals adminUserSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.svc.Update(r.Context(), command.UpdateUser{
		ID:         id,
		IsAdmin:    signals.IsAdmin,
		IsDisabled: signals.IsDisabled,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(h.renderUserList(users))
	sse.PatchElements(`<div id="modal"></div>`)
}

func (h *AdminHandler) renderUserList(users []entities.User) string {
	var sb strings.Builder
	sb.WriteString(`<div id="tab-content"><div class="level"><div class="level-left"><h2 class="title is-4">Users</h2></div></div>`)
	sb.WriteString(`<table class="table is-fullwidth is-striped is-hoverable">`)
	sb.WriteString(`<thead><tr><th>Email</th><th>Admin</th><th>Disabled</th><th>Created</th><th></th></tr></thead><tbody>`)

	for _, u := range users {
		admin := ""
		if u.IsAdmin {
			admin = "Yes"
		}
		disabled := ""
		if u.IsDisabled {
			disabled = "Yes"
		}
		sb.WriteString(fmt.Sprintf(
			`<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td><button class="button is-small is-info" data-on:click="@get('/admin/users/%s/edit')">Edit</button></td></tr>`,
			html.EscapeString(u.Email),
			admin,
			disabled,
			u.CreatedAt.Format("2006-01-02"),
			u.ID.String(),
		))
	}

	sb.WriteString(`</tbody></table></div>`)
	return sb.String()
}
