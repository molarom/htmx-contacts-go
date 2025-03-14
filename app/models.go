package app

import (
	"net/http"
	"net/url"
	"strconv"

	"gitlab.com/romalor/roxi"

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

type deletesParams struct {
	ids []int
}

func (c *deletesParams) Bind(data []byte) error {
	q, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	p := q["selected_contact_ids"]

	c.ids = make([]int, 0, len(p))
	for _, v := range p {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		c.ids = append(c.ids, int(i))
	}
	return nil
}

func parseDeletesParams(r *http.Request) (deletesParams, error) {
	p := deletesParams{}
	if err := roxi.Bind(r, &p); err != nil {
		return deletesParams{}, err
	}

	return p, nil
}
