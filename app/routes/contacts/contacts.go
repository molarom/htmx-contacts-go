package app

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/pkg/archiver"
	"gitlab.com/romalor/htmx-contacts/pkg/flash"
	"gitlab.com/romalor/htmx-contacts/pkg/middleware/htmx"
	"gitlab.com/romalor/htmx-contacts/pkg/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/pkg/tpl"
)

type handlers struct {
	tpls  *tpl.Bundle
	store *contacts.Store
}

func (h *handlers) Home(ctx context.Context, r *http.Request) error {
	return roxi.Redirect(ctx, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) List(ctx context.Context, r *http.Request) error {
	qp := r.URL.Query()
	search := qp.Get("q")
	page := "1"
	if p := qp.Get("page"); p == "" {
		page = "1"
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	var contacts contacts.Contacts
	if q := qp.Get("q"); q != "" {
		contacts = h.store.Search(q)
	} else {
		contacts = h.store.Page(p)
	}

	if htmx.Get(ctx).Trigger == "search" {
		return h.tpls.Render(roxi.GetWriter(ctx), "rows.html", tpl.Data{
			"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
			"search":   search,
			"contacts": contacts,
			"page":     p,
			"archiver": archiver.Default(),
		})
	}
	return h.tpls.Render(roxi.GetWriter(ctx), "index.html", tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"search":   search,
		"contacts": contacts,
		"page":     p,
		"archiver": archiver.Default(),
	})
}

func (h *handlers) Count(ctx context.Context, r *http.Request) error {
	return roxi.Respond(ctx,
		Resp("("+strconv.Itoa(h.store.Count())+" total Contacts)"),
	)
}

func (h *handlers) New(ctx context.Context, r *http.Request) error {
	return h.tpls.Render(roxi.GetWriter(ctx), "new.html", tpl.Data{
		"contact": contacts.Contact{},
		"errors":  nil,
	})
}

func (h *handlers) Create(ctx context.Context, r *http.Request) error {
	c, err := parseCreateForm(r)
	if err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "new.html", tpl.Data{
			"contact": contacts.Contact{},
			"errors":  err,
		})
	}

	if err := h.store.Create(c); err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "new.html", tpl.Data{
			"contact": c,
			"errors":  map[string]error{"Email": err},
		})
	}

	flash.Add(roxi.GetWriter(ctx), r, "Created New Contact!")
	return roxi.Redirect(ctx, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) Email(ctx context.Context, r *http.Request) error {
	err := h.store.Validate(contacts.Contact{
		Email: r.URL.Query().Get("email"),
	})
	if err != nil {
		return roxi.Respond(ctx, Resp(err.Error()))
	}
	return roxi.Respond(ctx, Resp{})
}

func (h *handlers) View(ctx context.Context, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	h.store.Get(int(id))
	return h.tpls.Render(roxi.GetWriter(ctx), "show.html", tpl.Data{
		"contact": h.store.Get(int(id)),
	})
}

func (h *handlers) Edit(ctx context.Context, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	return h.tpls.Render(roxi.GetWriter(ctx), "edit.html", tpl.Data{
		"contact": h.store.Get(int(id)),
	})
}

func (h *handlers) Update(ctx context.Context, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	uc, err := parseCreateForm(r)
	if err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "edit.html", tpl.Data{
			"contact": contacts.Contact{},
			"error":   err,
		})
	}
	uc.Id = int(id)

	_ = h.store.Update(uc)

	flash.Add(roxi.GetWriter(ctx), r, "Updated Contact!")
	return roxi.Redirect(ctx, r, "/contacts/"+r.PathValue("contact_id")+"/view", http.StatusMovedPermanently)
}

func (h *handlers) Delete(ctx context.Context, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	if htmx.Get(ctx).Trigger == "delete-btn" {
		if ok := h.store.Delete(int(id)); ok {
			flash.Add(roxi.GetWriter(ctx), r, "Deleted Contact!")
		}
		return roxi.Redirect(ctx, r, "/contacts/", http.StatusSeeOther)
	}
	return roxi.Respond(ctx, Resp(""))
}

func (h *handlers) Deletes(ctx context.Context, r *http.Request) error {
	qp, err := parseDeletesParams(r)
	if err != nil {
		return err
	}

	for _, id := range qp.ids {
		_ = h.store.Delete(id)
	}
	flash.Add(roxi.GetWriter(ctx), r, "Deleted Contacts!")

	return h.tpls.Render(roxi.GetWriter(ctx), "index.html", tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"contacts": h.store.Page(1),
		"page":     1,
	})
}
