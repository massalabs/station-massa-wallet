package types

import (
	"fmt"

	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"lukechampine.com/blake3"
)

const (
	AddressLastVersion = 0x00
)

// Address wraps the Object type to provide additional functionalities.
type Address struct {
	*object.Object
}

// NewAddressFromPublicKey creates a new Address from a given PublicKey.
func NewAddressFromPublicKey(pubKey *PublicKey) *Address {
	data := blake3.Sum256(pubKey.MustBytes())
	return &Address{Object: &object.Object{Kind: object.UserAddress, Version: AddressLastVersion, Data: data[:]}}
}

// validate ensures the Object.Kind is an address type and the version is supported.
func (a *Address) validate() error {
	err := a.Object.Validate(AddressLastVersion, object.UserAddress, object.SmartContractAddress)
	if err != nil {
		return err
	}

	return nil
}

// Custom YAML marshaller for Address
func (a Address) MarshalYAML() (interface{}, error) {
	data, err := a.MarshalText()

	return string(data), err
}

// Custom YAML unmarshaller for Address
func (a *Address) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}

	return a.UnmarshalText([]byte(data))
}

// MarshalText overloads the TextMarshaler interface for Address.
func (a *Address) MarshalText() ([]byte, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a.Object.MarshalText()
}

// String returns the Address as a string by marshaling it to text.
func (a *Address) String() (string, error) {
	data, err := a.MarshalText()
	if err != nil {
		return "", fmt.Errorf("failed to marshal address: %w", err)
	}

	return string(data), nil
}

// UnmarshalText overloads the TextUnmarshaler interface for Address.
func (a *Address) UnmarshalText(text []byte) error {
	if a.Object == nil {
		a.Object = &object.Object{
			Kind:    object.UserAddress,
			Version: 0x00,
			Data:    nil,
		}
	}

	if err := a.Object.UnmarshalText(text); err != nil {
		return err
	}

	return a.validate()
}

// MarshalBinary overloads the BinaryMarshaler interface for Address.
func (a *Address) MarshalBinary() ([]byte, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a.Object.MarshalBinary()
}

// UnmarshalBinary overloads the BinaryUnmarshaler interface for Address.
func (a *Address) UnmarshalBinary(data []byte) error {
	switch data[0] {
	case 0x00:
		a.Object.Kind = object.UserAddress
	case 0x01:
		a.Object.Kind = object.SmartContractAddress
	default:
		return object.ErrInvalidType
	}

	if err := a.Object.UnmarshalBinary(data[1:]); err != nil {
		return err
	}

	return a.validate()
}
