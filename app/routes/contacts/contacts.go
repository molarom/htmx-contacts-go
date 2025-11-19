package app

import (
	"context"
	"net/http"
	"strconv"

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
}

func (h *handlers) Home(ctx context.Context, r *http.Request) error {
	return roxi.Redirect(ctx, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) List(ctx context.Context, r *http.Request) error {
	var req listRequest
	if err := rika.Bind(r, &req); err != nil {
		return rika.BadRequest(err.Error())
	}

	var contacts contacts.Contacts
	if req.Search != "" {
		contacts = h.store.Search(req.Search)
	} else {
		contacts = h.store.Page(req.Page)
	}

	template := "index.html"
	if req.Trigger == "search" {
		template = "rows.html"
	}

	return h.tpls.Render(roxi.GetWriter(ctx), template, tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"search":   req.Search,
		"page":     req.Page,
		"contacts": contacts,
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
	var req createRequest
	if err := rika.Bind(r, &req); err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "new.html", tpl.Data{
			"contact": req,
			"errors": err,
		})
	}

	if err := h.store.Create(req.toContact()); err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "new.html", tpl.Data{
			"contact": req,
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
	var req updateRequest
	if err := rika.Bind(r, &req); err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "edit.html", tpl.Data{
			"contact": contacts.Contact{},
			"error":   err,
		})
	}

	_ = h.store.Update(req.toContact())

	flash.Add(roxi.GetWriter(ctx), r, "Updated Contact!")
	return roxi.Redirect(ctx, r, "/contacts/"+r.PathValue("contact_id")+"/view", http.StatusMovedPermanently)
}

func (h *handlers) Delete(ctx context.Context, r *http.Request) error {
	var req deleteRequest
	if err := rika.Bind(r, &req); err != nil {
		return err
	}

	if req.Trigger == "delete-btn" {
		if ok := h.store.Delete(req.ContactId); ok {
			flash.Add(roxi.GetWriter(ctx), r, "Deleted Contact!")
		}
		return roxi.Redirect(ctx, r, "/contacts/", http.StatusSeeOther)
	}
	return roxi.Respond(ctx, Resp(""))
}

func (h *handlers) Deletes(ctx context.Context, r *http.Request) error {
	var req deleteRequest
	if err := rika.Bind(r, &req); err != nil {
		return err
	}

	for _, id := range req.Ids {
		_ = h.store.Delete(id)
	}
	flash.Add(roxi.GetWriter(ctx), r, "Deleted Contacts!")

	return h.tpls.Render(roxi.GetWriter(ctx), "index.html", tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"contacts": h.store.Page(1),
		"page":     1,
	})
}
