// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// RequestType An enumeration of allowed request types.
//
// swagger:model RequestType
type RequestType string

func NewRequestType(value RequestType) *RequestType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated RequestType.
func (m RequestType) Pointer() *RequestType {
	return &m
}

const (

	// RequestTypeAddSignRule captures enum value "addSignRule"
	RequestTypeAddSignRule RequestType = "addSignRule"

	// RequestTypeRemoveSignRule captures enum value "removeSignRule"
	RequestTypeRemoveSignRule RequestType = "removeSignRule"
)

// for schema
var requestTypeEnum []interface{}

func init() {
	var res []RequestType
	if err := json.Unmarshal([]byte(`["addSignRule","removeSignRule"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		requestTypeEnum = append(requestTypeEnum, v)
	}
}

func (m RequestType) validateRequestTypeEnum(path, location string, value RequestType) error {
	if err := validate.EnumCase(path, location, value, requestTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this request type
func (m RequestType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateRequestTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this request type based on context it is used
func (m RequestType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}