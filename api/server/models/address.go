// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
)

// Address Account's address.
//
// swagger:model Address
type Address string

// Validate validates this address
func (m Address) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this address based on context it is used
func (m Address) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
