package types

import (
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/crypto"
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"github.com/stretchr/testify/assert"
)

func TestEncryptedPrivateKey_Marshal(t *testing.T) {
	data := []byte{2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83, 35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40, 28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246, 2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104, 64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135}
	// Functionally, this text below does not mean anything.
	text := "S1jgLNMTYLuaKMUVZVLVRfp9nMDvjCx8ov2FtyGQSf87oQej21kMDTi1tbzUmtFgFWMPhyepn2mXF9SMQqaAsN2PV8wB72eArCtCCeYCMupxNrM2yT78r"

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

		err := ad.UnmarshalBinary(append([]byte{0x00}, data...))
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
	secretBuffer := memguard.NewBufferFromBytes([]byte{216, 39, 16, 253, 102, 99, 172, 42, 205, 205, 17, 23, 123, 144, 171, 13, 91, 219, 194, 251, 186, 234, 11, 222, 23, 221, 6, 75, 22, 61, 235, 254, 45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147})
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
	})

	t.Run("PublicKey", func(t *testing.T) {
		samplePassword := memguard.NewBufferFromBytes([]byte("bonjour"))

		// Get the public key from the private key
		publicKey, err := sampleEncryptedPrivateKey.PublicKey(samplePassword, sampleSalt, sampleNonce)
		assert.NoError(t, err)

		expectedSignature := []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147}
		assert.Equal(t, expectedSignature, publicKey.Data)
	})
}
