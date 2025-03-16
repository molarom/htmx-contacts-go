package archive

import (
	"context"
	"io"
	"net/http"
	"os"

	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/archiver"
	"gitlab.com/romalor/htmx-contacts/flash"
	"gitlab.com/romalor/htmx-contacts/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/tpl"
)

type handlers struct {
	tpls  *tpl.Bundle
	store *contacts.Store
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
	f, err := os.Open(archiver.Default().File())
	if err != nil {
		return err
	}
	defer f.Close()

	w := roxi.GetWriter(ctx)
	w.Header().Set("Content-Disposition", "attachment; filename=archive.json")

	if _, err := io.Copy(w, f); err != nil {
		return err
	}

	return nil
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
