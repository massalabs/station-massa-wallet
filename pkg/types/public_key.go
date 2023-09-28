package types

import (
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
)

const (
	PublicKeyLastVersion = 0x00
)

type PublicKey struct {
	*object.Object
}

// validate ensures the Object.Kind is an PublicKey type and the version is supported.
func (a *PublicKey) validate() error {
	err := a.Object.Validate(PublicKeyLastVersion, object.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

// Custom YAML marshaller for PublicKey
func (pk PublicKey) MarshalYAML() (interface{}, error) {
	return pk.MarshalBinary()
}

// MarshalText overloads the TextMarshaler interface for PublicKey.
func (a *PublicKey) MarshalText() ([]byte, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a.Object.MarshalText()
}

// UnmarshalText overloads the TextUnmarshaler interface for PublicKey.
func (a *PublicKey) UnmarshalText(text []byte) error {
	if err := a.Object.UnmarshalText(text); err != nil {
		return err
	}

	return a.validate()
}

// MarshalBinary overloads the BinaryMarshaler interface for PublicKey.
func (a *PublicKey) MarshalBinary() ([]byte, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a.Object.MarshalBinary()
}

// UnmarshalBinary overloads the BinaryUnmarshaler interface for PublicKey.
func (a *PublicKey) UnmarshalBinary(data []byte) error {
	if err := a.Object.UnmarshalBinary(data); err != nil {
		return err
	}

	return a.validate()
}
