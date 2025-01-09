// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// RequestHandlerFunc turns a function with the right signature into a request handler
type RequestHandlerFunc func(RequestParams) middleware.Responder

// Handle executing the request and returning a response
func (fn RequestHandlerFunc) Handle(params RequestParams) middleware.Responder {
	return fn(params)
}

// RequestHandler interface for that can handle valid request params
type RequestHandler interface {
	Handle(RequestParams) middleware.Responder
}

// NewRequest creates a new http.Handler for the request operation
func NewRequest(ctx *middleware.Context, handler RequestHandler) *Request {
	return &Request{Context: ctx, Handler: handler}
}

/*
	Request swagger:route POST /api/accounts/{nickname}/request request

Request user approval through UI.
*/
type Request struct {
	Context *middleware.Context
	Handler RequestHandler
}

func (o *Request) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewRequestParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
