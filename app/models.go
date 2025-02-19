package app

import (
	"net/http"

	"gitlab.com/romalor/htmx-contacts/validator"
)

type Contact struct {
	Id     int    `json:"id"`
	First  string `json:"first" validate:"required"`
	Last   string `json:"last" validate:"required"`
	Phone  string `json:"phone" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Errors error  `json:"-"`
}

type Contacts []Contact

type ListPage struct {
	Search   string
	Contacts Contacts
}

func parseCreateForm(r *http.Request) (Contact, error) {
	c := Contact{
		First: r.FormValue("first_name"),
		Last:  r.FormValue("last_name"),
		Phone: r.FormValue("phone"),
		Email: r.FormValue("email"),
	}

	if err := validator.Verify(c); err != nil {
		return Contact{}, err
	}
	return c, nil
}
