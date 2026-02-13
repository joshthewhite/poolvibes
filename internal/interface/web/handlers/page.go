package handlers

import (
	"fmt"
	"html"
	"io/fs"
	"net/http"
	"strings"

	"github.com/joshthewhite/poolvibes/internal/application/services"
)

type PageHandler struct {
	layoutFS fs.FS
}

func NewPageHandler(layoutFS fs.FS) *PageHandler {
	return &PageHandler{layoutFS: layoutFS}
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(h.layoutFS, "templates/layout.html")
	if err != nil {
		http.Error(w, "failed to load layout", http.StatusInternalServerError)
		return
	}

	page := string(data)

	user, _ := services.UserFromContext(r.Context())
	if user != nil {
		navbarEnd := fmt.Sprintf(`<div class="navbar-end">
                    <span class="navbar-item has-text-light">%s</span>
                    <div class="navbar-item">
                        <form method="POST" action="/logout">
                            <button type="submit" class="button is-small is-light is-outlined">Logout</button>
                        </form>
                    </div>
                </div>`, html.EscapeString(user.Email))
		page = strings.Replace(page, `<div class="navbar-end">
                    <span class="navbar-item has-text-light">Pool Maintenance Manager</span>
                </div>`, navbarEnd, 1)

		if user.IsAdmin {
			adminTab := `                        <li data-class:is-active="$tab === 'admin'">
                            <a data-on:click="$tab = 'admin'; @get('/admin/users')">Admin</a>
                        </li>
                    </ul>`
			page = strings.Replace(page, `                    </ul>`, adminTab, 1)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(page))
}
