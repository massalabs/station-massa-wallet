package account

import (
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/types"
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"github.com/stretchr/testify/assert"
)

const (
	nickname       = "bonjour"
	password       = "bonjour"
	privateKeyText = "S12eCL2rGvRT4wZKaH7KdLd7fuhCF1Vt34SrNnRDEtduMZrjMxHz"
)

func TestNewAccount(t *testing.T) {
	// Create test values for the password and nickname
	samplePassword := memguard.NewBufferFromBytes([]byte(password))

	// Call the New function with the test values
	account, err := Generate(samplePassword, nickname)

	t.Run("ValidateAccountCreation", func(t *testing.T) {
		assert.NoError(t, err)

		assert.NotNil(t, account)
		assert.Equal(t, uint8(AccountLastVersion), *account.Version)
		assert.Equal(t, nickname, account.Nickname)
	})
}

func newAccount(t *testing.T) *Account {
	// Create test values for the password and nickname
	samplePassword := memguard.NewBufferFromBytes([]byte(password))
	privateKey := memguard.NewBufferFromBytes([]byte(privateKeyText))

	// Call the New function with the test values
	account, err := NewFromPrivateKey(samplePassword, nickname, privateKey)
	assert.NoError(t, err)

	return account
}

func TestNewAccountFromPrivateKey(t *testing.T) {
	account := newAccount(t)

	t.Run("ValidateAccountCreation", func(t *testing.T) {
		assert.NotNil(t, account)
		assert.Equal(t, uint8(AccountLastVersion), *account.Version)
		assert.Equal(t, nickname, account.Nickname)

		expectedPublicKey := []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147}
		assert.Equal(t, expectedPublicKey, account.PublicKey.Object.Data)

		expectedAddress := []byte{0x77, 0x13, 0x86, 0x8f, 0xe5, 0x5a, 0xd1, 0xdb, 0x9c, 0x8, 0x30, 0x7c, 0x61, 0x5e, 0xdf, 0xc0, 0xc8, 0x3b, 0x5b, 0xd9, 0x88, 0xec, 0x2e, 0x3c, 0xe9, 0xe4, 0x1c, 0xf1, 0xf9, 0x4d, 0xc5, 0xd1}
		assert.Equal(t, expectedAddress, account.Address.Object.Data)
	})
}

func TestPrivateKeyTextInClear(t *testing.T) {
	account := newAccount(t)

	samplePassword := memguard.NewBufferFromBytes([]byte(password))

	got, err := account.PrivateKeyTextInClear(samplePassword)
	assert.NoError(t, err)
	assert.Equal(t, privateKeyText, string(got.Bytes()))
	got.Destroy()
}

func TestSign(t *testing.T) {
	account := newAccount(t)
	samplePassword := memguard.NewBufferFromBytes([]byte(password))
	sampleData := []byte("Test")

	// Sign the data
	signature, err := account.Sign(samplePassword, sampleData)
	assert.NoError(t, err)

	// Assert
	assert.Greater(t, len(signature), 1)
	assert.Equal(t, byte(types.EncryptedPrivateKeyLastVersion), signature[0])
	expectedSignature := []byte{0x0, 0xe7, 0xeb, 0xd0, 0x39, 0xd3, 0xa3, 0x70, 0x70, 0xee, 0x38, 0xee, 0x95, 0x78, 0xd7, 0x3d, 0x7d, 0x74, 0xc4, 0x1a, 0x3, 0x1c, 0xfa, 0x3, 0xd4, 0x34, 0x1d, 0x67, 0x81, 0x64, 0x2c, 0xb7, 0xb0, 0x7c, 0xab, 0x30, 0xf1, 0x1d, 0x22, 0x39, 0x27, 0x7c, 0x9d, 0x5b, 0x4c, 0x9e, 0xcb, 0xa4, 0xe9, 0x8a, 0x5, 0x42, 0x20, 0xbb, 0x97, 0x7, 0x5e, 0x71, 0x87, 0x10, 0x40, 0xec, 0x8e, 0x62, 0x7}
	assert.Equal(t, expectedSignature, signature)

	t.Run("sign with new generated account", func(t *testing.T) {
		// Create test values for the password and nickname
		samplePassword := memguard.NewBufferFromBytes([]byte(password))

		// Call the New function with the test values
		account, err := Generate(samplePassword, nickname)
		assert.NoError(t, err)

		samplePassword = memguard.NewBufferFromBytes([]byte(password))

		signature, err := account.Sign(samplePassword, sampleData)
		assert.NoError(t, err)
		assert.Len(t, signature, 64+1) // +1 for the prefix
	})
}

func TestMarshal(t *testing.T) {
	accountText := `Version: 1
Nickname: bonjour
Address: AU1uSeKuYRQkC8fHWoohBmq8WKYxGTXuzaTcGomAoLXTGdq8kEsR
Salt: [145, 114, 211, 33, 247, 163, 215, 171, 90, 186, 97, 47, 43, 252, 68, 170]
Nonce: [113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63]
CipheredData: [2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83,
  35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40,
  28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246,
  2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104,
  64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135]
PublicKey: [0, 45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25,
  108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147]
`
	sampleSalt := [16]byte{145, 114, 211, 33, 247, 163, 215, 171, 90, 186, 97, 47, 43, 252, 68, 170}
	sampleNonce := [12]byte{113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63}
	sampleAccount, err := New(
		&[]uint8{AccountLastVersion}[0],
		nickname,
		&types.Address{
			Object: &object.Object{
				Kind:    object.UserAddress,
				Version: 0x00,
				Data:    []byte{0x77, 0x13, 0x86, 0x8f, 0xe5, 0x5a, 0xd1, 0xdb, 0x9c, 0x8, 0x30, 0x7c, 0x61, 0x5e, 0xdf, 0xc0, 0xc8, 0x3b, 0x5b, 0xd9, 0x88, 0xec, 0x2e, 0x3c, 0xe9, 0xe4, 0x1c, 0xf1, 0xf9, 0x4d, 0xc5, 0xd1},
			},
		},
		sampleSalt,
		sampleNonce,
		&types.EncryptedPrivateKey{
			Object: &object.Object{
				Kind:    object.EncryptedPrivateKey,
				Version: 0x00,
				Data:    []byte{2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83, 35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40, 28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246, 2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104, 64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135},
			},
		},
		&types.PublicKey{
			Object: &object.Object{
				Kind:    object.PublicKey,
				Version: 0x00,
				Data:    []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147},
			},
		},
	)
	assert.NoError(t, err)

	t.Run("Marshal", func(t *testing.T) {
		marshaledAccount, err := sampleAccount.Marshal()
		assert.NoError(t, err)

		assert.Equal(t, accountText, string(marshaledAccount))
	})

	t.Run("Unmarshal", func(t *testing.T) {
		account := NewEmpty()
		err := account.Unmarshal([]byte(accountText))
		assert.NoError(t, err)

		assert.Equal(t, uint8(1), *account.Version)
		assert.Equal(t, nickname, account.Nickname)
		assert.Equal(t, sampleSalt, account.Salt)
		assert.Equal(t, sampleNonce, account.Nonce)
		assert.Equal(t, sampleAccount.Address.Object.Data, account.Address.Object.Data)
		assert.Equal(t, sampleAccount.Address.Object.Version, account.Address.Object.Version)
		assert.Equal(t, sampleAccount.Address.Object.Kind, account.Address.Object.Kind)
		assert.Equal(t, sampleAccount.CipheredData.Object.Data, account.CipheredData.Object.Data)
		assert.Equal(t, sampleAccount.CipheredData.Object.Version, account.CipheredData.Object.Version)
		assert.Equal(t, sampleAccount.CipheredData.Object.Kind, account.CipheredData.Object.Kind)
		assert.Equal(t, sampleAccount.PublicKey.Object.Data, account.PublicKey.Object.Data)
		assert.Equal(t, sampleAccount.PublicKey.Object.Version, account.PublicKey.Object.Version)
		assert.Equal(t, sampleAccount.PublicKey.Object.Kind, account.PublicKey.Object.Kind)
	})

	t.Run("Unmarshal with missing field: version", func(t *testing.T) {
		accountText := `Nickname: bonjour
Address: AU1uSeKuYRQkC8fHWoohBmq8WKYxGTXuzaTcGomAoLXTGdq8kEsR
Salt: [145, 114, 211, 33, 247, 163, 215, 171, 90, 186, 97, 47, 43, 252, 68, 170]
Nonce: [113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63]
CipheredData: [2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83,
	35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40,
	28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246,
	2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104,
	64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135]
PublicKey: [0, 45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25,
	108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147]
`
		account := NewEmpty()
		err := account.Unmarshal([]byte(accountText))
		assert.ErrorContains(t, err, "missing version")
	})

	t.Run("Unmarshal with missing field: public key", func(t *testing.T) {
		accountText := `Version: 1
Nickname: bonjour
Address: AU1uSeKuYRQkC8fHWoohBmq8WKYxGTXuzaTcGomAoLXTGdq8kEsR
Salt: [145, 114, 211, 33, 247, 163, 215, 171, 90, 186, 97, 47, 43, 252, 68, 170]
Nonce: [113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63]
CipheredData: [2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83,
	35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40,
	28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246,
	2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104,
	64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135]
`
		account := NewEmpty()
		err := account.Unmarshal([]byte(accountText))
		assert.Error(t, err)
	})

	t.Run("Unmarshal with missing field: salt", func(t *testing.T) {
		accountText := `Version: 1
Nickname: bonjour
Address: AU1uSeKuYRQkC8fHWoohBmq8WKYxGTXuzaTcGomAoLXTGdq8kEsR
Nonce: [113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63]
CipheredData: [2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83,
	35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40,
	28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246,
	2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104,
	64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135]
PublicKey: [0, 45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25,
	108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147]
`
		account := NewEmpty()
		err := account.Unmarshal([]byte(accountText))
		assert.Error(t, err)
	})
}

func BenchmarkProtectWallet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		iteration()
	}
}

func iteration() {
	samplePassword := memguard.NewBufferFromBytes([]byte("password"))

	_, err := Generate(samplePassword, "nickname")
	if err != nil {
		panic(err)
	}
}
