package handler

import (
	"errors"
	"html/template"
	"net/http"

	"robo-ruka/internal/domain"
	"robo-ruka/internal/service"
)

type Handler struct {
	svc  *service.Service
	tmpl *template.Template
}

func New(svc *service.Service, tmpl *template.Template) *Handler {
	return &Handler{svc: svc, tmpl: tmpl}
}

type viewModel struct {
	IsOn  bool
	Error string
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var (
		st     domain.Status
		errMsg string
	)

	if raw := r.URL.Query().Get("status"); raw != "" {
		updated, err := h.svc.Status.Update(raw)
		switch {
		case errors.Is(err, service.ErrInvalidStatus):
			errMsg = "ожидается status=on или status=off"
		case err != nil:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		default:
			st = updated
		}
	}

	if st == "" {
		cur, err := h.svc.Status.Current()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		st = cur
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.Execute(w, viewModel{IsOn: st.IsOn(), Error: errMsg}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
