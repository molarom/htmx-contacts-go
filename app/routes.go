package app

import (
	"gitlab.com/romalor/roxi"

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

	// Homepage
	mux.GET("/", h.Home)
	mux.GET("/contacts", h.List)

	// Archive
	mux.GET("/contacts/archive", h.Status)
	mux.POST("/contacts/archive", h.Archive)
	mux.GET("/contacts/archive/file", h.ArchiveFile)

	// Contacts
	mux.GET("/contacts/count", h.Count)
	mux.GET("/contacts/new", h.New)
	mux.POST("/contacts/new", h.Create)
	mux.GET("/contacts/email", h.Email)
	mux.GET("/contacts/view/:contact_id", h.View)

	mux.GET("/contacts/:contact_id/edit", h.Edit)
	mux.POST("/contacts/:contact_id/edit", h.Update)

	mux.DELETE("/contacts/:contact_id", h.Delete)
	mux.DELETE("/contacts", h.Deletes)
}
