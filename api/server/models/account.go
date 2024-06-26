// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Account Account object (V0).
//
// swagger:model Account
type Account struct {

	// address
	// Required: true
	Address Address `json:"address"`

	// balance
	// Required: true
	Balance Amount `json:"balance"`

	// candidate balance
	// Required: true
	CandidateBalance Amount `json:"candidateBalance"`

	// key pair
	// Required: true
	KeyPair KeyPair `json:"keyPair"`

	// nickname
	// Required: true
	Nickname Nickname `json:"nickname"`

	// status
	// Enum: ["ok","corrupted"]
	Status string `json:"status,omitempty"`
}

// Validate validates this account
func (m *Account) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBalance(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCandidateBalance(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKeyPair(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNickname(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Account) validateAddress(formats strfmt.Registry) error {

	if err := validate.Required("address", "body", Address(m.Address)); err != nil {
		return err
	}

	if err := m.Address.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("address")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("address")
		}
		return err
	}

	return nil
}

func (m *Account) validateBalance(formats strfmt.Registry) error {

	if err := validate.Required("balance", "body", Amount(m.Balance)); err != nil {
		return err
	}

	if err := m.Balance.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("balance")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("balance")
		}
		return err
	}

	return nil
}

func (m *Account) validateCandidateBalance(formats strfmt.Registry) error {

	if err := validate.Required("candidateBalance", "body", Amount(m.CandidateBalance)); err != nil {
		return err
	}

	if err := m.CandidateBalance.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("candidateBalance")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("candidateBalance")
		}
		return err
	}

	return nil
}

func (m *Account) validateKeyPair(formats strfmt.Registry) error {

	if err := m.KeyPair.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("keyPair")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("keyPair")
		}
		return err
	}

	return nil
}

func (m *Account) validateNickname(formats strfmt.Registry) error {

	if err := validate.Required("nickname", "body", Nickname(m.Nickname)); err != nil {
		return err
	}

	if err := m.Nickname.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("nickname")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("nickname")
		}
		return err
	}

	return nil
}

var accountTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["ok","corrupted"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		accountTypeStatusPropEnum = append(accountTypeStatusPropEnum, v)
	}
}

const (

	// AccountStatusOk captures enum value "ok"
	AccountStatusOk string = "ok"

	// AccountStatusCorrupted captures enum value "corrupted"
	AccountStatusCorrupted string = "corrupted"
)

// prop value enum
func (m *Account) validateStatusEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, accountTypeStatusPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Account) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this account based on the context it is used
func (m *Account) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAddress(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateBalance(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCandidateBalance(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateKeyPair(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNickname(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Account) contextValidateAddress(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Address.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("address")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("address")
		}
		return err
	}

	return nil
}

func (m *Account) contextValidateBalance(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Balance.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("balance")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("balance")
		}
		return err
	}

	return nil
}

func (m *Account) contextValidateCandidateBalance(ctx context.Context, formats strfmt.Registry) error {

	if err := m.CandidateBalance.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("candidateBalance")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("candidateBalance")
		}
		return err
	}

	return nil
}

func (m *Account) contextValidateKeyPair(ctx context.Context, formats strfmt.Registry) error {

	if err := m.KeyPair.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("keyPair")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("keyPair")
		}
		return err
	}

	return nil
}

func (m *Account) contextValidateNickname(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Nickname.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("nickname")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("nickname")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Account) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Account) UnmarshalBinary(b []byte) error {
	var res Account
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
