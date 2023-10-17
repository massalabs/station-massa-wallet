package crypto_test

import (
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

// isWiped checks if all elements in the byte array are 0x00.
func isWiped(b []byte) bool {
	for _, v := range b {
		if v != 0x00 {
			return false
		}
	}

	return true
}

func TestSealAndUnsealSecret(t *testing.T) {
	tests := []struct {
		name     string
		password []byte
		salt     [16]byte
		nonce    [12]byte
		secret   []byte
	}{
		{
			name:     "Test Case 1",
			password: []byte("password1"),
			salt:     [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			nonce:    [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			secret:   []byte("secret1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize cipher
			aeadCipher, secretKey, err := crypto.NewSecretCipher(tt.password, tt.salt[:])
			assert.NoError(t, err, "NewSecretCipher should not return an error")

			// Seal the secret
			secretBuffer := memguard.NewBufferFromBytes(tt.secret)
			encryptedSecret := crypto.SealSecret(aeadCipher, tt.nonce[:], secretBuffer)

			// Check if the tt.secret array has been wiped to 0x00
			assert.True(t, isWiped(tt.secret), "Secret should be wiped")

			// Unseal the secret
			decryptedBuffer, err := crypto.UnsealSecret(aeadCipher, tt.nonce[:], encryptedSecret)
			assert.NoError(t, err, "UnsealSecret should not return an error")

			// Destroy the secret Key buffer as we don't need aeadCipher anymore
			secretKey.Destroy()

			// Check if the decrypted secret matches the original
			assert.Equal(t, string(secretBuffer.Bytes()), string(decryptedBuffer.Bytes()), "Decrypted secret should match the original")

			// Destroy the decrypted buffer
			decryptedBuffer.Destroy()
		})
	}
}
