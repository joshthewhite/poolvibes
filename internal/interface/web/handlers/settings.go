package handlers

import (
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
	"github.com/starfederation/datastar-go/datastar"
)

type SettingsHandler struct {
	svc *services.UserService
}

func NewSettingsHandler(svc *services.UserService) *SettingsHandler {
	return &SettingsHandler{svc: svc}
}

type settingsSignals struct {
	Phone       string `json:"settingsPhone"`
	NotifyEmail bool   `json:"settingsNotifyEmail"`
	NotifySMS   bool   `json:"settingsNotifySms"`
}

func (h *SettingsHandler) Page(w http.ResponseWriter, r *http.Request) {
	user, err := services.UserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.SettingsPage(user.Phone, user.NotifyEmail, user.NotifySMS))
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var signals settingsSignals
	if err := datastar.ReadSignals(r, &signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.svc.UpdatePreferences(r.Context(), command.UpdateNotificationPreferences{
		Phone:       signals.Phone,
		NotifyEmail: signals.NotifyEmail,
		NotifySMS:   signals.NotifySMS,
	})
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.PatchElementTempl(templates.SettingsMessage("is-danger is-light", err.Error()))
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(templates.SettingsMessage("is-success is-light", "Settings saved successfully."))
}
