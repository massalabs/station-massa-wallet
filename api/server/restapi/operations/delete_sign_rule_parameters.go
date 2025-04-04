// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewDeleteSignRuleParams creates a new DeleteSignRuleParams object
//
// There are no default values defined in the spec.
func NewDeleteSignRuleParams() DeleteSignRuleParams {

	return DeleteSignRuleParams{}
}

// DeleteSignRuleParams contains all the bound params for the delete sign rule operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteSignRule
type DeleteSignRuleParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Account's short name.
	  Required: true
	  In: path
	*/
	Nickname string
	/*The ID of the sign rule
	  Required: true
	  In: path
	*/
	RuleID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteSignRuleParams() beforehand.
func (o *DeleteSignRuleParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rNickname, rhkNickname, _ := route.Params.GetOK("nickname")
	if err := o.bindNickname(rNickname, rhkNickname, route.Formats); err != nil {
		res = append(res, err)
	}

	rRuleID, rhkRuleID, _ := route.Params.GetOK("ruleId")
	if err := o.bindRuleID(rRuleID, rhkRuleID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindNickname binds and validates parameter Nickname from path.
func (o *DeleteSignRuleParams) bindNickname(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.Nickname = raw

	return nil
}

// bindRuleID binds and validates parameter RuleID from path.
func (o *DeleteSignRuleParams) bindRuleID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.RuleID = raw

	return nil
}
