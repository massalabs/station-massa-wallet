package types

import (
	"crypto/ed25519"
	"fmt"

	"github.com/awnumar/memguard"
	"github.com/btcsuite/btcutil/base58"
	"github.com/massalabs/station-massa-wallet/pkg/crypto"
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"lukechampine.com/blake3"
)

const (
	// EncryptedPrivateKeyLastVersion is the last version of the encrypted private key.
	EncryptedPrivateKeyLastVersion = 0x00
)

type EncryptedPrivateKey struct {
	*object.Object
}

// validate ensures the Object.Kind is an EncryptedPrivateKey type and the version is supported.
func (e *EncryptedPrivateKey) validate() error {
	err := e.Object.Validate(EncryptedPrivateKeyLastVersion, object.EncryptedPrivateKey)
	if err != nil {
		return err
	}

	return nil
}

// Custom YAML marshaller for EncryptedPrivateKey
func (e EncryptedPrivateKey) MarshalYAML() (interface{}, error) {
	return e.MarshalBinary()
}

// Custom YAML unmarshaller for EncryptedPrivateKey
func (e *EncryptedPrivateKey) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data []byte
	if err := unmarshal(&data); err != nil {
		return err
	}

	return e.UnmarshalBinary(data)
}

// MarshalText overloads the TextMarshaler interface for EncryptedPrivateKey.
func (e *EncryptedPrivateKey) MarshalText() ([]byte, error) {
	if err := e.validate(); err != nil {
		return nil, err
	}

	return e.Object.MarshalText()
}

// UnmarshalText overloads the TextUnmarshaler interface for EncryptedPrivateKey.
func (e *EncryptedPrivateKey) UnmarshalText(text []byte) error {
	if err := e.Object.UnmarshalText(text); err != nil {
		return err
	}

	return e.validate()
}

// MarshalBinary overloads the BinaryMarshaler interface for EncryptedPrivateKey.
func (e *EncryptedPrivateKey) MarshalBinary() ([]byte, error) {
	if err := e.validate(); err != nil {
		return nil, err
	}

	return e.Object.MarshalBinary()
}

// UnmarshalBinary overloads the BinaryUnmarshaler interface for EncryptedPrivateKey.
func (e *EncryptedPrivateKey) UnmarshalBinary(data []byte) error {
	if e.Object == nil {
		e.Object = &object.Object{
			Kind:    object.EncryptedPrivateKey,
			Version: 0x00,
			Data:    nil,
		}
	}

	e.Object.Version = 0x00 // We can't know the version at this stage. The version is in the plain text. `data` is ciphered.
	e.Object.Data = data

	return e.validate()
}

// Sign signs the given data using the private key. Password is destroyed.
func (e *EncryptedPrivateKey) Sign(password *memguard.LockedBuffer, salt, nonce, data []byte) ([]byte, error) {
	digest := blake3.Sum256(data)

	privateKeyInClear, err := privateKey(password, salt, nonce, e.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	defer privateKeyInClear.Destroy()

	return append([]byte{EncryptedPrivateKeyLastVersion}, ed25519.Sign(privateKeyInClear.Bytes(), digest[:])...), nil
}

// SignWithPrivateKey signs the given data using the private key. Private key is destroyed.
func (e *EncryptedPrivateKey) SignWithPrivateKey(privateKey *memguard.LockedBuffer, data []byte) []byte {
	digest := blake3.Sum256(data)

	defer privateKey.Destroy()

	return append([]byte{EncryptedPrivateKeyLastVersion}, ed25519.Sign(privateKey.Bytes(), digest[:])...)
}

// PublicKey returns the public key corresponding to the private key. Password is destroyed.
func (e *EncryptedPrivateKey) PublicKey(password *memguard.LockedBuffer, salt, nonce []byte) (*PublicKey, error) {
	privateKeyInClear, err := privateKey(password, salt, nonce, e.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	publicKeyBytes := ed25519.PrivateKey(privateKeyInClear.Bytes()).Public().(ed25519.PublicKey)

	privateKeyInClear.Destroy()

	return &PublicKey{
		Object: &object.Object{
			Data:    publicKeyBytes,
			Kind:    object.PublicKey,
			Version: PublicKeyLastVersion,
		},
	}, nil
}

// PrivateKeyTextInClear returns the private key in clear. Password is destroyed.
func (e *EncryptedPrivateKey) PrivateKeyTextInClear(password *memguard.LockedBuffer, salt, nonce []byte) (*memguard.LockedBuffer, error) {
	privateKeyInClear, err := e.PrivateKeyBytesInClear(password, salt, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	seed := ed25519.PrivateKey(privateKeyInClear.Bytes()).Seed()
	privateKeyInClear.Destroy()

	seedBuffer := memguard.NewBufferFromBytes(seed)

	privateKey := e.Kind.Prefix() + base58.CheckEncode(seedBuffer.Bytes(), e.Version)
	seedBuffer.Destroy()
	privateKeyBuffer := memguard.NewBufferFromBytes([]byte(privateKey))

	return privateKeyBuffer, nil
}

func (e *EncryptedPrivateKey) PrivateKeyBytesInClear(password *memguard.LockedBuffer, salt, nonce []byte) (*memguard.LockedBuffer, error) {
	privateKeyInClear, err := privateKey(password, salt, nonce, e.Data)
	if err != nil {
		return nil, fmt.Errorf("PrivateKeyTextInClear: %w", err)
	}

	return privateKeyInClear, nil
}

// HasAccess returns true if the password is valid for the account. It destroys the password.
func (e *EncryptedPrivateKey) HasAccess(password *memguard.LockedBuffer, salt, nonce []byte) bool {
	privateKeyInClear, err := privateKey(password, salt, nonce, e.Data)
	if err != nil {
		return false
	}

	privateKeyInClear.Destroy()

	return true
}

// privateKey returns the private key in clear. Password is destroyed.
func privateKey(password *memguard.LockedBuffer, salt, nonce, encryptedKey []byte) (*memguard.LockedBuffer, error) {
	aeadCipher, secretKey, err := crypto.NewSecretCipher(password.Bytes(), salt[:])
	defer password.Destroy()
	defer secretKey.Destroy()

	if err != nil {
		return nil, err
	}

	secret, err := crypto.UnsealSecret(aeadCipher, nonce[:], encryptedKey)
	if err != nil {
		return nil, err
	}

	data := secret.Bytes()
	result := memguard.NewBuffer(64)
	result.Copy(data[1:])
	secret.Destroy()

	return result, nil
}
