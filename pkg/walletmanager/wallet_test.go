package walletmanager

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/massalabs/station-massa-wallet/pkg/types"
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	// Clean
	clean(t)

	var w *Wallet
	sampleSalt := [16]byte{145, 114, 211, 33, 247, 163, 215, 171, 90, 186, 97, 47, 43, 252, 68, 170}
	sampleNonce := [12]byte{113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63}
	sampleNickname := "bonjour2"
	sampleAccount, err := account.New(
		&[]uint8{account.AccountLastVersion}[0],
		sampleNickname,
		types.Address{
			Object: &object.Object{
				Kind:    object.UserAddress,
				Version: types.AddressLastVersion,
				Data:    []byte{0x77, 0x13, 0x86, 0x8f, 0xe5, 0x5a, 0xd1, 0xdb, 0x9c, 0x8, 0x30, 0x7c, 0x61, 0x5e, 0xdf, 0xc0, 0xc8, 0x3b, 0x5b, 0xd9, 0x88, 0xec, 0x2e, 0x3c, 0xe9, 0xe4, 0x1c, 0xf1, 0xf9, 0x4d, 0xc5, 0xd1},
			},
		},
		sampleSalt,
		sampleNonce,
		types.EncryptedPrivateKey{
			Object: &object.Object{
				Kind:    object.EncryptedPrivateKey,
				Version: types.EncryptedPrivateKeyLastVersion,
				Data:    []byte{2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83, 35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40, 28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246, 2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104, 64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135},
			},
		},
		types.PublicKey{
			Object: &object.Object{
				Kind:    object.PublicKey,
				Version: types.PublicKeyLastVersion,
				Data:    []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147},
			},
		},
	)
	assert.NoError(t, err)

	t.Run("Create Wallet", func(t *testing.T) {
		newWallet := New()
		w = newWallet
		assert.NotNil(t, w)
	})

	t.Run("Add Account", func(t *testing.T) {
		err := w.AddAccount(sampleAccount)
		assert.NoError(t, err)

		assert.Equal(t, 1, w.GetAccountCount())
		accountPath, err := w.accountPath(sampleNickname)
		assert.NoError(t, err)
		assert.FileExists(t, accountPath)
	})

	t.Run("Add Account: nickname not unique", func(t *testing.T) {
		err := w.AddAccount(sampleAccount)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "nickname is not unique: this account name already exists")

		assert.Equal(t, 1, w.GetAccountCount())
	})

	t.Run("Add Account: address not unique", func(t *testing.T) {
		sampleAccount, err := account.New(
			&[]uint8{account.AccountLastVersion}[0],
			"bonjour3",
			types.Address{
				Object: &object.Object{
					Kind:    object.UserAddress,
					Version: types.AddressLastVersion,
					Data:    []byte{0x77, 0x13, 0x86, 0x8f, 0xe5, 0x5a, 0xd1, 0xdb, 0x9c, 0x8, 0x30, 0x7c, 0x61, 0x5e, 0xdf, 0xc0, 0xc8, 0x3b, 0x5b, 0xd9, 0x88, 0xec, 0x2e, 0x3c, 0xe9, 0xe4, 0x1c, 0xf1, 0xf9, 0x4d, 0xc5, 0xd1},
				},
			},
			sampleSalt,
			sampleNonce,
			types.EncryptedPrivateKey{
				Object: &object.Object{
					Kind:    object.EncryptedPrivateKey,
					Version: types.EncryptedPrivateKeyLastVersion,
					Data:    []byte{2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83, 35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40, 28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246, 2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104, 64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135},
				},
			},
			types.PublicKey{
				Object: &object.Object{
					Kind:    object.PublicKey,
					Version: types.PublicKeyLastVersion,
					Data:    []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147},
				},
			},
		)
		assert.NoError(t, err)

		err = w.AddAccount(sampleAccount)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "address is not unique: this account address already exists")

		assert.Equal(t, 1, w.GetAccountCount())

		sampleAccount.Nickname = "bonjour2"
	})

	t.Run("Get Account", func(t *testing.T) {
		acc, err := w.GetAccount(sampleNickname)
		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, uint8(1), *acc.Version)
		assert.Equal(t, sampleNickname, acc.Nickname)
		assert.Equal(t, sampleSalt, acc.Salt)
		assert.Equal(t, sampleNonce, acc.Nonce)
		assert.Equal(t, sampleAccount.Address.Object.Data, acc.Address.Object.Data)
		assert.Equal(t, sampleAccount.Address.Object.Version, acc.Address.Object.Version)
		assert.Equal(t, sampleAccount.Address.Object.Kind, acc.Address.Object.Kind)
		assert.Equal(t, sampleAccount.CipheredData.Object.Data, acc.CipheredData.Object.Data)
		assert.Equal(t, sampleAccount.CipheredData.Object.Version, acc.CipheredData.Object.Version)
		assert.Equal(t, sampleAccount.CipheredData.Object.Kind, acc.CipheredData.Object.Kind)
		assert.Equal(t, sampleAccount.PublicKey.Object.Data, acc.PublicKey.Object.Data)
		assert.Equal(t, sampleAccount.PublicKey.Object.Version, acc.PublicKey.Object.Version)
		assert.Equal(t, sampleAccount.PublicKey.Object.Kind, acc.PublicKey.Object.Kind)
	})

	t.Run("Get Account: not found", func(t *testing.T) {
		acc, err := w.GetAccount("wrong nickname")
		assert.Error(t, err)
		assert.Nil(t, acc)
	})

	t.Run("Get Account: new file added manually", func(t *testing.T) {
		// User can add an account file in the account folder by its own,
		// we want to wallet to be able to manage this account.
		nickname := "unit-test"

		accountPath, err := w.accountPath(nickname)
		assert.NoError(t, err)

		copy(t, "../../tests/wallet_unit-test.yaml", accountPath)

		acc, err := w.GetAccount(nickname)
		assert.NoError(t, err)
		assert.Equal(t, uint8(1), *acc.Version)
		assert.Equal(t, nickname, acc.Nickname)
		assert.Equal(t, 2, w.GetAccountCount())
	})

	t.Run("Delete Account", func(t *testing.T) {
		err := w.DeleteAccount(sampleNickname)
		assert.NoError(t, err)
		assert.Equal(t, 1, w.GetAccountCount())

		accountPath, err := w.accountPath(sampleNickname)
		assert.NoError(t, err)
		assert.NoFileExists(t, accountPath)
	})

	t.Run("Delete Account: not found", func(t *testing.T) {
		err := w.DeleteAccount("unknown")
		assert.Error(t, err)
		assert.Equal(t, 1, w.GetAccountCount())
	})

	// Clean
	clean(t)
}

func clean(t *testing.T) {
	accountsPath, err := Path()
	assert.NoError(t, err)

	files, err := os.ReadDir(accountsPath)
	assert.NoError(t, err)

	for _, f := range files {
		fileName := f.Name()
		filePath := path.Join(accountsPath, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			os.Remove(filePath)
		}
	}

	assert.NoError(t, err)
}

func copy(t *testing.T, src string, dst string) {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	assert.NoError(t, err)
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	assert.NoError(t, err)
}
