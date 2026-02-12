package handlers

import (
	"io/fs"
	"net/http"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}
