package app

import (
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
	mux.GET("/", h.Home)
	mux.GET("/contacts", h.List)
	mux.GET("/contacts/new", h.New)
	mux.POST("/contacts/new", h.Create)
	mux.GET("/contacts/view/:contact_id", h.View)
	mux.GET("/contacts/:contact_id/edit", h.Edit)
	mux.POST("/contacts/:contact_id/edit", h.Update)
	mux.DELETE("/contacts/:contact_id/delete", h.Delete)
}
