// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetAllAssetsParams creates a new GetAllAssetsParams object
//
// There are no default values defined in the spec.
func NewGetAllAssetsParams() GetAllAssetsParams {

	return GetAllAssetsParams{}
}

// GetAllAssetsParams contains all the bound params for the get all assets operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetAllAssets
type GetAllAssetsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*List of asset addresses (MRC-20 token addresses) to retrieve info for.
	  Required: true
	  In: query
	*/
	AssetAddresses []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetAllAssetsParams() beforehand.
func (o *GetAllAssetsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qAssetAddresses, qhkAssetAddresses, _ := qs.GetOK("assetAddresses")
	if err := o.bindAssetAddresses(qAssetAddresses, qhkAssetAddresses, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAssetAddresses binds and validates array parameter AssetAddresses from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetAllAssetsParams) bindAssetAddresses(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("assetAddresses", "query", rawData)
	}
	var qvAssetAddresses string
	if len(rawData) > 0 {
		qvAssetAddresses = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	assetAddressesIC := swag.SplitByFormat(qvAssetAddresses, "")
	if len(assetAddressesIC) == 0 {
		return errors.Required("assetAddresses", "query", assetAddressesIC)
	}

	var assetAddressesIR []string
	for _, assetAddressesIV := range assetAddressesIC {
		assetAddressesI := assetAddressesIV

		assetAddressesIR = append(assetAddressesIR, assetAddressesI)
	}

	o.AssetAddresses = assetAddressesIR

	return nil
}