package handlers

import (
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	email := ""
	isAdmin := false
	user, _ := services.UserFromContext(r.Context())
	if user != nil {
		email = user.Email
		isAdmin = user.IsAdmin
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.Layout(email, isAdmin).Render(r.Context(), w)
}
