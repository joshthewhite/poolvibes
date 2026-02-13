package handlers

import (
	"fmt"
	"html"
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
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

	notifyEmail := "false"
	if user.NotifyEmail {
		notifyEmail = "true"
	}
	notifySMS := "false"
	if user.NotifySMS {
		notifySMS = "true"
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`<div id="tab-content">
		<div data-signals:settingsPhone="'%s'" data-signals:settingsNotifyEmail="%s" data-signals:settingsNotifySms="%s">
			<div class="level">
				<div class="level-left"><h2 class="title is-4">Notification Settings</h2></div>
			</div>
			<div id="settings-message"></div>
			<div class="box" style="max-width: 500px;">
				<div class="field">
					<label class="label">Phone Number</label>
					<div class="control">
						<input data-bind:settingsPhone type="tel" class="input" placeholder="+15551234567">
					</div>
					<p class="help">Required for SMS notifications. Include country code.</p>
				</div>
				<div class="field">
					<label class="checkbox">
						<input data-bind:settingsNotifyEmail type="checkbox"> Email notifications
					</label>
					<p class="help">Receive email alerts when tasks are due.</p>
				</div>
				<div class="field">
					<label class="checkbox">
						<input data-bind:settingsNotifySms type="checkbox"> SMS notifications
					</label>
					<p class="help">Receive text message alerts when tasks are due.</p>
				</div>
				<div class="field mt-4">
					<div class="control">
						<button class="button is-link" data-on:click="@put('/settings')">Save Settings</button>
					</div>
				</div>
			</div>
		</div>
	</div>`,
		html.EscapeString(user.Phone),
		notifyEmail,
		notifySMS,
	))
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
		sse.PatchElements(fmt.Sprintf(`<div id="settings-message"><div class="notification is-danger is-light">%s</div></div>`, html.EscapeString(err.Error())))
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.PatchElements(`<div id="settings-message"><div class="notification is-success is-light">Settings saved successfully.</div></div>`)
}
