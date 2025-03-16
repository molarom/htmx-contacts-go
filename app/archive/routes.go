package archive

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

	// Archive
	mux.GET("/contacts/archive", h.Status)
	mux.POST("/contacts/archive", h.Archive)
	mux.GET("/contacts/archive/file", h.ArchiveFile)
}
