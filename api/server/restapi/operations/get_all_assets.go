// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/massalabs/station-massa-wallet/api/server/models"
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
	GetAllAssets swagger:route GET /api/assets getAllAssets

Retrieve all assets (MRC-20) info.
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

// GetAllAssetsOKBody get all assets o k body
//
// swagger:model GetAllAssetsOKBody
type GetAllAssetsOKBody struct {

	// assets
	Assets []models.AssetInfo `json:"assets"`
}

// Validate validates this get all assets o k body
func (o *GetAllAssetsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAssets(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetAllAssetsOKBody) validateAssets(formats strfmt.Registry) error {
	if swag.IsZero(o.Assets) { // not required
		return nil
	}

	for i := 0; i < len(o.Assets); i++ {

		if err := o.Assets[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getAllAssetsOK" + "." + "assets" + "." + strconv.Itoa(i))
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getAllAssetsOK" + "." + "assets" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

// ContextValidate validate this get all assets o k body based on the context it is used
func (o *GetAllAssetsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateAssets(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetAllAssetsOKBody) contextValidateAssets(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Assets); i++ {

		if err := o.Assets[i].ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getAllAssetsOK" + "." + "assets" + "." + strconv.Itoa(i))
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getAllAssetsOK" + "." + "assets" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetAllAssetsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetAllAssetsOKBody) UnmarshalBinary(b []byte) error {
	var res GetAllAssetsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
