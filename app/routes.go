package app

import (
	"github.com/gorilla/sessions"
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
		sessions.NewCookieStore([]byte("securestring")),
	}
	mux.Handle("GET", "/", h.Home)
	mux.Handle("GET", "/contacts", h.List)
	mux.Handle("GET", "/contacts/new", h.New)
	mux.Handle("GET", "/contacts/view/:contact_id", h.View)
	mux.Handle("POST", "/contacts/new", h.Create)
}
