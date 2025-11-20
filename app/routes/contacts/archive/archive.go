package archive

import (
	"context"
	"net/http"

	"gitlab.com/romalor/rika"
	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/pkg/archiver"
	"gitlab.com/romalor/htmx-contacts/pkg/flash"
	"gitlab.com/romalor/htmx-contacts/pkg/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/pkg/tpl"
)

type handlers struct {
	tpls  *tpl.Bundle
	store *contacts.Store
	fr *rika.FileResponder
}

func (h *handlers) Archive(ctx context.Context, r *http.Request) error {
	go archiver.Default().Run()
	return h.tpls.Render(roxi.GetWriter(ctx), "archive_ui.html", tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"contacts": h.store.Page(1),
		"page":     1,
		"archiver": archiver.Default(),
	})
}

func (h *handlers) Status(ctx context.Context, r *http.Request) error {
	return h.tpls.Render(roxi.GetWriter(ctx), "archive_ui.html", tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"contacts": h.store.Page(1),
		"page":     1,
		"archiver": archiver.Default(),
	})
}

func (h *handlers) ArchiveFile(ctx context.Context, r *http.Request) error {
	return h.fr.Attachment(ctx, r, archiver.Default().File())
}

func (h *handlers) Reset(ctx context.Context, r *http.Request) error {
	archiver.Default().Reset()
	return h.tpls.Render(roxi.GetWriter(ctx), "archive_ui.html", tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"contacts": h.store.Page(1),
		"page":     1,
		"archiver": archiver.Default(),
	})
}
