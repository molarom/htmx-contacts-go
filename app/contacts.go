package app

import (
	"net/http"
	"strconv"

	"gitlab.com/romalor/htmx-contacts/tpl"
)

type handlers struct {
	tpls         *tpl.Bundle
	contactStore Contacts
}

func (h *handlers) Home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) List(w http.ResponseWriter, r *http.Request) {
	h.tpls.Render(w, "index.html", ListPage{
		r.URL.Query().Get("q"),
		h.contactStore,
	})
}

func (h *handlers) New(w http.ResponseWriter, r *http.Request) {
	h.tpls.Render(w, "new.html", Contact{})
}

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {
	c, err := parseCreateForm(r)
	if err != nil {
		h.tpls.Render(w, "new.html", Contact{Errors: err})
		return
	}

	c.Id = len(h.contactStore) + 1
	h.contactStore = append(h.contactStore, c)

	http.Redirect(w, r, "/contacts", http.StatusMovedPermanently)
}

func (h *handlers) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("contact_id"), 10, 64)
	if err != nil {
		return
	}

	contact := Contact{}
	for _, c := range h.contactStore {
		if c.Id == int(id) {
			contact = c
		}
	}

	h.tpls.Render(w, "show.html", contact)
}
