// Package object provides functionalities for handling different types of Massa objects.
// It supports encoding and decoding operations with appropriate prefixes and versioning.
package object

import (
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

// Kind represents the type of data being manipulated.
type Kind int

const (
	Unknown Kind = iota
	EncryptedPrivateKey
	PublicKey
	UserAddress
	SmartContractAddress
)

// Prefix returns the prefix corresponding to the kind of data.
func (k Kind) Prefix() string {
	switch k {
	case EncryptedPrivateKey:
		return "S"
	case PublicKey:
		return "P"
	case UserAddress:
		return "AU"
	case SmartContractAddress:
		return "AS"
	default:
		return ""
	}
}

// NewKind determines the kind of data based on its prefix.
func NewKind(data string) Kind {
	if len(data) < 2 {
		return Unknown
	}

	switch data[0] {
	case 'S':
		return EncryptedPrivateKey
	case 'P':
		return PublicKey
	case 'A':
		if data[1] == 'U' {
			return UserAddress
		} else if data[1] == 'S' {
			return SmartContractAddress
		}
	}

	return Unknown
}

var (
	ErrUnsupportedVersion = errors.New("unsupported version")
	ErrEmptyType          = errors.New("empty type")
	ErrInvalidType        = errors.New("invalid type")
)

// Object represents a Massa object with its data, kind, and version.
type Object struct {
	Data    []byte
	Kind    Kind
	Version byte
}

// Validate ensures the Object.Kind is one of the expected types and the version is supported.
func (o *Object) Validate(lastVersion byte, expectedKinds ...Kind) error {
	if o.Version > lastVersion {
		return ErrUnsupportedVersion
	}

	for _, kind := range expectedKinds {
		if o.Kind == kind {
			return nil
		}
	}

	return ErrInvalidType
}

// MarshalText implements the encoding.TextMarshaler interface.
func (o *Object) MarshalText() ([]byte, error) {
	dataStr := base58.CheckEncode(o.Data, o.Version)
	return []byte(o.Kind.Prefix() + dataStr), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (o *Object) UnmarshalText(text []byte) error {
	o.Kind = NewKind(string(text))
	var err error

	switch o.Kind {
	case EncryptedPrivateKey, PublicKey:
		o.Data, o.Version, err = base58.CheckDecode(string(text[1:]))
	case UserAddress, SmartContractAddress:
		o.Data, o.Version, err = base58.CheckDecode(string(text[2:]))
	}

	return err
}

// MustBytes returns the binary representation of the object.
func (o *Object) MustBytes() []byte {
	data, err := o.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return data
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (o *Object) MarshalBinary() ([]byte, error) {
	switch o.Kind {
	case UserAddress:
		return append([]byte{0x00, o.Version}, o.Data...), nil
	case SmartContractAddress:
		return append([]byte{0x01, o.Version}, o.Data...), nil
	default:
		return append([]byte{o.Version}, o.Data...), nil
	}
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (o *Object) UnmarshalBinary(data []byte) error {
	o.Version = data[0]
	o.Data = data[1:]

	return nil
}
