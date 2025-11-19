package app

import (
	"gitlab.com/romalor/htmx-contacts/pkg/stores/contacts"
	"gitlab.com/romalor/htmx-contacts/pkg/validator"
)

// ----------------------------------------------------------------------
// Requests

type listRequest struct {
	// Headers
	Trigger string `header:"hx-trigger"`

	// Query Params
	Page int `query:"page"`
	Search string `query:"q"`

}

func (r *listRequest) Validate() error {
	if r.Page < 1 {
		r.Page = 1
	}
	return nil
}

type createRequest struct {
	First string `form:"first_name" validate:"required"`
	Last  string `form:"last_name" validate:"required"`
	Phone string `form:"phone" validate:"required"`
	Email string `form:"email" validate:"required,email"`
}

func (r *createRequest) Validate() error {
	return validator.Verify(r)
}

func (r *createRequest) toContact() contacts.Contact {
	return contacts.Contact{
		First: r.First,
		Last: r.Last,
		Phone: r.Phone,
		Email: r.Email,
	}
}

type updateRequest struct {
	Id int `path:"contact_id"`
	createRequest
}

func (r *updateRequest) toContact() contacts.Contact {
	return contacts.Contact{
		Id: r.Id,
		First: r.First,
		Last: r.Last,
		Phone: r.Phone,
		Email: r.Email,
	}
}

type deleteRequest struct {
	Trigger string `header:"hx-trigger"`
	ContactId int `path:"contact_id"`
	Ids []int `query:"selected_contact_ids"`
}


// ----------------------------------------------------------------------
// Responses

type Resp []byte

func (r Resp) Response() ([]byte, string, error) {
	return r, "text/plain; charset=utf-8", nil
}

