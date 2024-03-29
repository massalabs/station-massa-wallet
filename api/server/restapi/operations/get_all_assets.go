// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetAllAssetsHandlerFunc turns a function with the right signature into a get all assets handler
type GetAllAssetsHandlerFunc func(GetAllAssetsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAllAssetsHandlerFunc) Handle(params GetAllAssetsParams) middleware.Responder {
	return fn(params)
}

// GetAllAssetsHandler interface for that can handle valid get all assets params
type GetAllAssetsHandler interface {
	Handle(GetAllAssetsParams) middleware.Responder
}

// NewGetAllAssets creates a new http.Handler for the get all assets operation
func NewGetAllAssets(ctx *middleware.Context, handler GetAllAssetsHandler) *GetAllAssets {
	return &GetAllAssets{Context: ctx, Handler: handler}
}

/*
	GetAllAssets swagger:route GET /api/accounts/{nickname}/assets getAllAssets

Get all assets with their balance.
*/
type GetAllAssets struct {
	Context *middleware.Context
	Handler GetAllAssetsHandler
}

func (o *GetAllAssets) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetAllAssetsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
