package types

import (
	"crypto/ed25519"

	"github.com/massalabs/station-massa-wallet/pkg/types/object"
)

const (
	PublicKeyLastVersion = 0x00
)

type PublicKey struct {
	*object.Object
}

// validate ensures the Object.Kind is an PublicKey type and the version is supported.
func (p *PublicKey) validate() error {
	if len(p.Data) != ed25519.PublicKeySize {
		return object.ErrUnsupportedVersion
	}

	err := p.Object.Validate(PublicKeyLastVersion, object.PublicKey)
	if err != nil {
		return err
	}

	return nil
}

// Custom YAML marshaller for PublicKey
func (p PublicKey) MarshalYAML() (interface{}, error) {
	return p.MarshalBinary()
}

// Custom YAML unmarshaller for PublicKey
func (pk *PublicKey) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data []byte
	if err := unmarshal(&data); err != nil {
		return err
	}

	return pk.UnmarshalBinary(data)
}

// MarshalText overloads the TextMarshaler interface for PublicKey.
func (p *PublicKey) MarshalText() ([]byte, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	return p.Object.MarshalText()
}

// UnmarshalText overloads the TextUnmarshaler interface for PublicKey.
func (p *PublicKey) UnmarshalText(text []byte) error {
	if err := p.Object.UnmarshalText(text); err != nil {
		return err
	}

	return p.validate()
}

// MarshalBinary overloads the BinaryMarshaler interface for PublicKey.
func (p *PublicKey) MarshalBinary() ([]byte, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	return p.Object.MarshalBinary()
}

// UnmarshalBinary overloads the BinaryUnmarshaler interface for PublicKey.
func (p *PublicKey) UnmarshalBinary(data []byte) error {
	if p.Object == nil {
		p.Object = &object.Object{
			Kind:    object.PublicKey,
			Version: 0x00,
			Data:    nil,
		}
	}

	if err := p.Object.UnmarshalBinary(data); err != nil {
		return err
	}

	return p.validate()
}

func (p *PublicKey) VerifySignature(data, signature []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(p.Object.Data), data, signature)
}
