package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"

	"github.com/awnumar/memguard"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// SecretKeySizeInBytes specifies the length of the secret key in bytes.
	// This size is defined according to the Massa standard.
	SecretKeySizeInBytes = 32

	// Pkdf2NbRound is the number of rounds for the PBKDF2 algorithm.
	// This number of rounds is defined according to the Massa standard.
	Pkdf2NbRound = 600_000

	// SaltSizeInBytes is the size of the salt in bytes.
	SaltSizeInBytes = 16
	// NonceSizeInBytes is the size of the nonce in bytes.
	NonceSizeInBytes = 12
)

// NewSecretCipher initializes a new AEAD cipher using AES-GCM and PBKDF2 for key derivation.
// Note: The returned locked buffer containing the secret key is to be destroy as soon as possible.
func NewSecretCipher(password []byte, salt []byte) (cipher.AEAD, *memguard.LockedBuffer, error) {
	secretKey := pbkdf2.Key(password, salt, Pkdf2NbRound, SecretKeySizeInBytes, sha256.New)

	// this also wipes the original secretKey slice
	guardedKey := memguard.NewBufferFromBytes(secretKey)

	block, err := aes.NewCipher(guardedKey.Bytes())
	if err != nil {
		guardedKey.Destroy()
		return nil, nil, fmt.Errorf("initializing AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		guardedKey.Destroy()
		return nil, nil, fmt.Errorf("initializing AES-GCM: %w", err)
	}

	return aesGCM, guardedKey, nil
}

// UnsealSecret decrypts an encrypted secret using the provided AEAD cipher and nonce.
// The decrypted secret is securely stored in a memguard locked buffer.
// Note: The returned locked buffer is to be destroy as soon as possible.
func UnsealSecret(aeadCipher cipher.AEAD, nonce []byte, encryptedSecret []byte) (*memguard.LockedBuffer, error) {
	secret, err := aeadCipher.Open(nil, nonce, encryptedSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// this also wipes the original secret slice
	guardedSecret := memguard.NewBufferFromBytes(secret)

	return guardedSecret, nil
}

// SealSecret encrypts a given secret using the provided AEAD cipher and nonce.
func SealSecret(aeadCipher cipher.AEAD, nonce []byte, secret *memguard.LockedBuffer) []byte {
	return aeadCipher.Seal(nil, nonce, secret.Bytes(), nil)
}
