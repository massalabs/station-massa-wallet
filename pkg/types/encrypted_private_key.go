package types

import (
	"crypto/ed25519"

	"github.com/awnumar/memguard"
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
func (a *EncryptedPrivateKey) validate() error {
	err := a.Object.Validate(EncryptedPrivateKeyLastVersion, object.EncryptedPrivateKey)
	if err != nil {
		return err
	}

	return nil
}

// MarshalText overloads the TextMarshaler interface for EncryptedPrivateKey.
func (a *EncryptedPrivateKey) MarshalText() ([]byte, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a.Object.MarshalText()
}

// UnmarshalText overloads the TextUnmarshaler interface for EncryptedPrivateKey.
func (a *EncryptedPrivateKey) UnmarshalText(text []byte) error {
	if err := a.Object.UnmarshalText(text); err != nil {
		return err
	}
	return a.validate()
}

// MarshalBinary overloads the BinaryMarshaler interface for EncryptedPrivateKey.
func (a *EncryptedPrivateKey) MarshalBinary() ([]byte, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a.Object.MarshalBinary()
}

// UnmarshalBinary overloads the BinaryUnmarshaler interface for EncryptedPrivateKey.
func (a *EncryptedPrivateKey) UnmarshalBinary(data []byte) error {
	if err := a.Object.UnmarshalBinary(data); err != nil {
		return err
	}
	return a.validate()
}

// Sign signs the given data using the private key.
// password is destroyed just after being used even if an error occurs.
func (e *EncryptedPrivateKey) Sign(password *memguard.LockedBuffer, salt, nonce, data []byte) ([]byte, error) {
	digest := blake3.Sum256(data)

	privateKeyInClear, err := PrivateKey(password, salt, nonce, e.Data)
	if err != nil {
		return nil, err
	}

	defer privateKeyInClear.Destroy()

	return append([]byte{EncryptedPrivateKeyLastVersion}, ed25519.Sign(privateKeyInClear.Bytes(), digest[:])...), nil
}

// PublicKey returns the public key corresponding to the private key.
func (e *EncryptedPrivateKey) PublicKey(password *memguard.LockedBuffer, salt, nonce []byte) (*PublicKey, error) {
	privateKeyInClear, err := PrivateKey(password, salt, nonce, e.Data)
	if err != nil {
		return nil, err
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

func PrivateKey(password *memguard.LockedBuffer, salt, nonce, encryptedKey []byte) (*memguard.LockedBuffer, error) {
	aeadCipher, secretKey, err := crypto.NewSecretCipher(password.Bytes(), salt[:])
	defer password.Destroy()
	defer secretKey.Destroy()

	if err != nil {
		return nil, err
	}

	return crypto.UnsealSecret(aeadCipher, nonce[:], encryptedKey)
}
