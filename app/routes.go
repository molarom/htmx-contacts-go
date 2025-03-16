package app

import (
	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/app/archive"
	"gitlab.com/romalor/htmx-contacts/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/tpl"
)

type Config struct {
	TplBundle *tpl.Bundle
	Store     *contacts.Store
}

func Routes(mux *roxi.Mux, cfg Config) {
	h := &handlers{
		cfg.TplBundle,
		cfg.Store,
	}

	// Archive
	archive.Routes(mux, archive.Config{
		TplBundle: cfg.TplBundle,
		Store:     cfg.Store,
	})

	// Homepage
	mux.GET("/", h.Home)
	mux.GET("/contacts", h.List)
	mux.GET("/contacts/count", h.Count)

	// Contacts
	mux.GET("/contacts/new", h.New)
	mux.POST("/contacts/new", h.Create)
	mux.GET("/contacts/email", h.Email)
	mux.GET("/contacts/:contact_id/view", h.View)

	// Edits
	mux.GET("/contacts/:contact_id/edit", h.Edit)
	mux.POST("/contacts/:contact_id/edit", h.Update)

	// Deletes
	mux.DELETE("/contacts/:contact_id", h.Delete)
	mux.DELETE("/contacts", h.Deletes)
}
