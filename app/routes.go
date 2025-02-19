package app

import (
	"net/http"

	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/tpl"
)

type Config struct {
	TplBundle    *tpl.Bundle
	ContactStore Contacts
}

func Routes(mux *roxi.Mux, cfg Config) {
	h := &handlers{
		cfg.TplBundle,
		cfg.ContactStore,
	}
	mux.HandlerFunc("GET", "/", h.Home)
	mux.HandlerFunc("GET", "/contacts", h.List)
	mux.HandlerFunc("GET", "/contacts/new", h.New)
	mux.HandlerFunc("GET", "/contacts/view/:contact_id", h.View)
	mux.HandlerFunc("POST", "/contacts/new", h.Create)
}

func StdRoutes(mux *http.ServeMux, cfg Config) {
	h := &handlers{
		cfg.TplBundle,
		cfg.ContactStore,
	}
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /contacts", h.List)
	mux.HandleFunc("GET /contacts/new", h.New)
	mux.HandleFunc("GET /contacts/view/:contact_id", h.View)
	mux.HandleFunc("POST /contacts/new", h.Create)
}
