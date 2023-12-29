package types

import (
	"crypto/ed25519"
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/crypto"
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"github.com/stretchr/testify/assert"
	"lukechampine.com/blake3"
)

func TestEncryptedPrivateKey_Marshal(t *testing.T) {
	data := []byte{0x2, 0x56, 0x85, 0x92, 0x52, 0xb8, 0xc1, 0xa0, 0x78, 0x2c, 0xc6, 0xd1, 0x45, 0xe6, 0x53, 0x23, 0x24, 0xeb, 0x12, 0x69, 0x4a, 0x75, 0xe4, 0xed, 0x70, 0x41, 0x20, 0x0, 0xfa, 0xb4, 0xc7, 0x1a}
	// Functionally, this text below does not mean anything.
	text := "S122inbQk89WNigayDBb2PhuNN2rmhdr5BUVN4tzVoN3kkZePQn"

	t.Run("Marshal Text", func(t *testing.T) {
		tests := []struct {
			name         string
			privateKey   EncryptedPrivateKey
			expectedData []byte
			expectedErr  error
		}{
			{
				name: "Valid private key",
				privateKey: EncryptedPrivateKey{
					Object: &object.Object{
						Kind:    object.EncryptedPrivateKey,
						Version: 0x00,
						Data:    data,
					},
				},
				expectedData: []byte(text),
				expectedErr:  nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := tt.privateKey.MarshalText()
				assert.Equal(t, err, tt.expectedErr)
				assert.Equal(t, string(tt.expectedData), string(got))
			})
		}
	})

	t.Run("Unmarshal Text", func(t *testing.T) {
		ad := EncryptedPrivateKey{
			Object: &object.Object{
				Kind:    object.EncryptedPrivateKey,
				Version: 0x00,
				Data:    data,
			},
		}

		err := ad.UnmarshalText([]byte(text))
		assert.NoError(t, err)

		assert.Equal(t, data, ad.Object.Data)
	})

	t.Run("Marshal Binary", func(t *testing.T) {
		ad := EncryptedPrivateKey{
			Object: &object.Object{
				Kind:    object.EncryptedPrivateKey,
				Version: 0x00,
				Data:    data,
			},
		}

		b, err := ad.MarshalBinary()
		assert.NoError(t, err)
		assert.Equal(t, data, b)
	})

	t.Run("UnMarshal Binary", func(t *testing.T) {
		ad := EncryptedPrivateKey{
			Object: &object.Object{
				Kind:    object.EncryptedPrivateKey,
				Version: 0x00,
				Data:    data,
			},
		}

		err := ad.UnmarshalBinary(data)
		assert.NoError(t, err)
		assert.Equal(t, data, ad.Object.Data)
	})
}

func TestEncryptedPrivateKey(t *testing.T) {
	// Declare sample data
	sampleData := []byte("Test")
	sampleSalt := []byte{145, 114, 211, 33, 247, 163, 215, 171, 90, 186, 97, 47, 43, 252, 68, 170}
	sampleNonce := []byte{113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63}
	samplePassword := memguard.NewBufferFromBytes([]byte("bonjour"))

	// Prepare the secret
	aeadCipher, secretKey, err := crypto.NewSecretCipher(samplePassword.Bytes(), sampleSalt)
	assert.NoError(t, err)
	// secretBuffer is the secret key in clear without the version.
	secretBuffer := memguard.NewBufferFromBytes([]byte{0, 216, 39, 16, 253, 102, 99, 172, 42, 205, 205, 17, 23, 123, 144, 171, 13, 91, 219, 194, 251, 186, 234, 11, 222, 23, 221, 6, 75, 22, 61, 235, 254, 45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147})
	encryptedSecret := crypto.SealSecret(aeadCipher, sampleNonce, secretBuffer)

	secretKey.Destroy()

	// Prepare the encrypted private key object
	sampleEncryptedPrivateKey := &EncryptedPrivateKey{
		Object: &object.Object{
			Kind:    object.EncryptedPrivateKey,
			Version: 0x00,
			Data:    encryptedSecret,
		},
	}

	t.Run("Sign", func(t *testing.T) {
		// Sign the data
		signature, err := sampleEncryptedPrivateKey.Sign(samplePassword, sampleSalt, sampleNonce, sampleData)
		assert.NoError(t, err)

		// Assert
		assert.Greater(t, len(signature), 1)
		assert.Equal(t, byte(EncryptedPrivateKeyLastVersion), signature[0])
		expectedSignature := []byte{0x0, 0xe7, 0xeb, 0xd0, 0x39, 0xd3, 0xa3, 0x70, 0x70, 0xee, 0x38, 0xee, 0x95, 0x78, 0xd7, 0x3d, 0x7d, 0x74, 0xc4, 0x1a, 0x3, 0x1c, 0xfa, 0x3, 0xd4, 0x34, 0x1d, 0x67, 0x81, 0x64, 0x2c, 0xb7, 0xb0, 0x7c, 0xab, 0x30, 0xf1, 0x1d, 0x22, 0x39, 0x27, 0x7c, 0x9d, 0x5b, 0x4c, 0x9e, 0xcb, 0xa4, 0xe9, 0x8a, 0x5, 0x42, 0x20, 0xbb, 0x97, 0x7, 0x5e, 0x71, 0x87, 0x10, 0x40, 0xec, 0x8e, 0x62, 0x7}
		assert.Equal(t, expectedSignature, signature)

		// Get the public key and verify signature
		samplePassword = memguard.NewBufferFromBytes([]byte("bonjour")) // recreate a memguard buffer from the password
		publicKey, err := sampleEncryptedPrivateKey.PublicKey(samplePassword, sampleSalt, sampleNonce)
		assert.NoError(t, err)

		digest := blake3.Sum256(sampleData)
		assert.True(t, ed25519.Verify(publicKey.Data, digest[:], expectedSignature[1:]))
	})

	t.Run("PublicKey", func(t *testing.T) {
		samplePassword := memguard.NewBufferFromBytes([]byte("bonjour"))

		// Get the public key from the private key
		publicKey, err := sampleEncryptedPrivateKey.PublicKey(samplePassword, sampleSalt, sampleNonce)
		assert.NoError(t, err)

		expectedPublicKey := []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147}
		assert.Equal(t, expectedPublicKey, publicKey.Data)
	})

	t.Run("PrivateKeyBytesInClear", func(t *testing.T) {
		samplePassword := memguard.NewBufferFromBytes([]byte("bonjour"))

		// Get the public key from the private key
		privateKeyBytesInClear, err := sampleEncryptedPrivateKey.PrivateKeyBytesInClear(samplePassword, sampleSalt, sampleNonce)
		assert.NoError(t, err)

		expectedPrivateKeyBytesInClear := []byte{0xd8, 0x27, 0x10, 0xfd, 0x66, 0x63, 0xac, 0x2a, 0xcd, 0xcd, 0x11, 0x17, 0x7b, 0x90, 0xab, 0xd, 0x5b, 0xdb, 0xc2, 0xfb, 0xba, 0xea, 0xb, 0xde, 0x17, 0xdd, 0x6, 0x4b, 0x16, 0x3d, 0xeb, 0xfe, 0x2d, 0x96, 0xbc, 0xda, 0xcb, 0xbe, 0x41, 0x38, 0x2c, 0xa2, 0x3e, 0x52, 0xe3, 0xd2, 0x19, 0x6c, 0xba, 0x65, 0xe7, 0xa1, 0xac, 0xd2, 0x9, 0xdf, 0xc9, 0x5c, 0x6b, 0x32, 0xb6, 0xa1, 0x8a, 0x93}
		assert.Equal(t, expectedPrivateKeyBytesInClear, privateKeyBytesInClear.Bytes())
	})
}
