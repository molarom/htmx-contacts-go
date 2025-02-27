package app

import (
	"net/http"

	"gitlab.com/romalor/htmx-contacts/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/validator"
)

func parseCreateForm(r *http.Request) (contacts.Contact, error) {
	c := contacts.Contact{
		First: r.FormValue("first_name"),
		Last:  r.FormValue("last_name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
	}

	if err := validator.Verify(c); err != nil {
		return contacts.Contact{}, err
	}
	return c, nil
}

type Resp []byte

func (r Resp) Response() ([]byte, string, error) {
	return r, "text/plain; charset=utf-8", nil
}
