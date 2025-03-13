package app

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/flash"
	"gitlab.com/romalor/htmx-contacts/middleware/htmx"
	"gitlab.com/romalor/htmx-contacts/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/tpl"
)

type handlers struct {
	tpls  *tpl.Bundle
	store *contacts.Store
}

func (h *handlers) Home(ctx context.Context, r *http.Request) error {
	return roxi.Redirect(ctx, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) List(ctx context.Context, r *http.Request) error {
	pg := r.URL.Query().Get("page")
	if pg == "" {
		pg = "1"
	}

	p, err := strconv.Atoi(pg)
	if err != nil {
		return err
	}

	var contacts contacts.Contacts
	if q := r.URL.Query().Get("q"); q != "" {
		contacts = h.store.Search(q)
		if htmx.Get(ctx).Trigger == "search" {
			return h.tpls.Render(roxi.GetWriter(ctx), "rows.html", tpl.Data{
				"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
				"contacts": contacts,
				"page":     p,
			})
		}
	} else {
		contacts = h.store.Page(p)
	}
	return h.tpls.Render(roxi.GetWriter(ctx), "index.html", tpl.Data{
		"flashes":  flash.Messages(roxi.GetWriter(ctx), r),
		"contacts": contacts,
		"page":     p,
	})
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

	h.store.Update(uc)

	flash.Add(roxi.GetWriter(ctx), r, "Updated Contact!")
	return roxi.Redirect(ctx, r, "/contacts/view/"+r.PathValue("contact_id"), http.StatusMovedPermanently)
}

func (h *handlers) Delete(ctx context.Context, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	if ok := h.store.Delete(int(id)); ok {
		flash.Add(roxi.GetWriter(ctx), r, "Deleted Contact!")
	}

	return roxi.Redirect(ctx, r, "/contacts/", http.StatusSeeOther)
}
