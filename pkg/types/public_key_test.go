package types

import (
	"testing"

	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"github.com/stretchr/testify/assert"
)

func TestAddress_MarshalText(t *testing.T) {
	data := []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147}
	text := "P1M5WJVWdjSWZ5xndJKXAnKFU9E8CUSbdXJRiuAexBpuDsc4QPK"

	t.Run("Marshal Text", func(t *testing.T) {
		tests := []struct {
			name         string
			publicKey    PublicKey
			expectedData []byte
			expectedErr  error
		}{
			{
				name: "Valid public key",
				publicKey: PublicKey{
					Object: &object.Object{
						Kind:    object.PublicKey,
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
				got, err := tt.publicKey.MarshalText()
				assert.Equal(t, err, tt.expectedErr)
				assert.Equal(t, string(tt.expectedData), string(got))
			})
		}
	})

	t.Run("Unmarshal Text", func(t *testing.T) {
		ad := PublicKey{
			Object: &object.Object{
				Kind:    object.PublicKey,
				Version: 0x00,
				Data:    data,
			},
		}

		err := ad.UnmarshalText([]byte(text))
		assert.NoError(t, err)

		assert.Equal(t, data, ad.Object.Data)
	})

	t.Run("Marshal Binary", func(t *testing.T) {
		ad := PublicKey{
			Object: &object.Object{
				Kind:    object.PublicKey,
				Version: 0x00,
				Data:    data,
			},
		}

		b, err := ad.MarshalBinary()
		assert.NoError(t, err)
		// we expect the version to be prepended to the data
		assert.Equal(t, append([]byte{0x00}, data...), b)
	})

	t.Run("UnMarshal Binary", func(t *testing.T) {
		ad := PublicKey{
			Object: &object.Object{
				Kind:    object.PublicKey,
				Version: 0x00,
				Data:    data,
			},
		}

		err := ad.UnmarshalBinary(append([]byte{0x00}, data...))
		assert.NoError(t, err)
		assert.Equal(t, data, ad.Object.Data)
	})
}
