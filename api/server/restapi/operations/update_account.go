// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateAccountHandlerFunc turns a function with the right signature into a update account handler
type UpdateAccountHandlerFunc func(UpdateAccountParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateAccountHandlerFunc) Handle(params UpdateAccountParams) middleware.Responder {
	return fn(params)
}

// UpdateAccountHandler interface for that can handle valid update account params
type UpdateAccountHandler interface {
	Handle(UpdateAccountParams) middleware.Responder
}

// NewUpdateAccount creates a new http.Handler for the update account operation
func NewUpdateAccount(ctx *middleware.Context, handler UpdateAccountHandler) *UpdateAccount {
	return &UpdateAccount{Context: ctx, Handler: handler}
}

/*
	UpdateAccount swagger:route PUT /api/accounts/{nickname} updateAccount

Update the account associated with the provided nickname in the path.
*/
type UpdateAccount struct {
	Context *middleware.Context
	Handler UpdateAccountHandler
}

func (o *UpdateAccount) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateAccountParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
