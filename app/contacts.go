package app

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/tpl"
)

type handlers struct {
	tpls         *tpl.Bundle
	contactStore Contacts
	session      sessions.Store
}

func (h *handlers) Home(ctx context.Context, r *http.Request) error {
	session, err := h.session.New(r, strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		return err
	}

	session.Save(r, roxi.GetWriter(ctx))
	return roxi.Redirect(ctx, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) List(ctx context.Context, r *http.Request) error {
	session, err := h.session.Get(r, strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		return err
	}

	return h.tpls.Render(roxi.GetWriter(ctx), "index.html", tpl.Data{
		"contacts": h.contactStore,
		"flashes":  session.Flashes(),
	})
}

func (h *handlers) New(ctx context.Context, r *http.Request) error {
	return h.tpls.Render(roxi.GetWriter(ctx), "new.html", tpl.Data{
		"contact": Contact{},
		"errors":  nil,
	})
}

func (h *handlers) Create(ctx context.Context, r *http.Request) error {
	session, err := h.session.Get(r, strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		return err
	}

	c, err := parseCreateForm(r)
	if err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "new.html", tpl.Data{
			"contact": Contact{},
			"errors":  err,
		})
	}

	c.Id = len(h.contactStore) + 1
	h.contactStore = append(h.contactStore, c)

	session.AddFlash("Created New Contact!")
	session.Save(r, roxi.GetWriter(ctx))
	return roxi.Redirect(ctx, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) View(ctx context.Context, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	for _, c := range h.contactStore {
		if c.Id != int(id) {
			continue
		}
		return h.tpls.Render(roxi.GetWriter(ctx), "show.html", tpl.Data{
			"contact": c,
		})
	}
	return h.tpls.Render(roxi.GetWriter(ctx), "show.html", tpl.Data{
		"contact": Contact{},
	})
}

func (h *handlers) Edit(ctx context.Context, r *http.Request) error {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	for _, c := range h.contactStore {
		if c.Id != int(id) {
			continue
		}
		return h.tpls.Render(roxi.GetWriter(ctx), "edit.html", tpl.Data{
			"contact": c,
		})
	}
	return h.tpls.Render(roxi.GetWriter(ctx), "edit.html", tpl.Data{
		"contact": Contact{},
	})
}

func (h *handlers) Update(ctx context.Context, r *http.Request) error {
	session, err := h.session.Get(r, strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return err
	}

	uc, err := parseCreateForm(r)
	if err != nil {
		return h.tpls.Render(roxi.GetWriter(ctx), "edit.html", tpl.Data{
			"contact": Contact{},
			"errorl":  err,
		})
	}
	uc.Id = int(id)

	for i, c := range h.contactStore {
		if c.Id != int(id) {
			continue
		}
		h.contactStore[i] = uc
		break
	}

	session.AddFlash(err)
	session.Save(r, roxi.GetWriter(ctx))
	return roxi.Redirect(ctx, r, "/contacts/view/"+r.PathValue("contact_id"), http.StatusMovedPermanently)
}
