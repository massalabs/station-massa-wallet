package wallet

import (
	"log"
	"os"
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/types"
	"github.com/massalabs/station-massa-wallet/pkg/types/object"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	if err := logger.InitializeGlobal("./unit-test.log"); err != nil {
		log.Fatalf("while initializing global logger: %s", err.Error())
	}

	walletPath, err := os.MkdirTemp(os.TempDir(), "*-wallet-dir")
	assert.NoError(t, err)

	var w *Wallet
	samplePassword := memguard.NewBufferFromBytes([]byte("password"))
	sampleSalt := [16]byte{145, 114, 211, 33, 247, 163, 215, 171, 90, 186, 97, 47, 43, 252, 68, 170}
	sampleNonce := [12]byte{113, 122, 168, 123, 48, 187, 178, 12, 209, 91, 243, 63}
	sampleAddressData := []byte{0x77, 0x13, 0x86, 0x8f, 0xe5, 0x5a, 0xd1, 0xdb, 0x9c, 0x8, 0x30, 0x7c, 0x61, 0x5e, 0xdf, 0xc0, 0xc8, 0x3b, 0x5b, 0xd9, 0x88, 0xec, 0x2e, 0x3c, 0xe9, 0xe4, 0x1c, 0xf1, 0xf9, 0x4d, 0xc5, 0xd1}
	sampleAddressText := "AU1uSeKuYRQkC8fHWoohBmq8WKYxGTXuzaTcGomAoLXTGdq8kEsR"
	addressOnlyRequiredFields := "AU1AAQExqUbw2PvBjNdZgodNg9jUFwxAfPP8mASPbamK3unxmXtm"
	samplePrivateKeyData := []byte{2, 86, 133, 146, 82, 184, 193, 160, 120, 44, 198, 209, 69, 230, 83, 35, 36, 235, 18, 105, 74, 117, 228, 237, 112, 65, 32, 0, 250, 180, 199, 26, 40, 28, 76, 116, 162, 95, 0, 103, 172, 8, 41, 11, 240, 185, 188, 215, 56, 170, 246, 2, 14, 16, 27, 214, 137, 103, 89, 111, 85, 149, 191, 38, 2, 43, 8, 183, 149, 104, 64, 149, 10, 106, 102, 156, 242, 178, 254, 189, 135}
	samplePublicKeyData := []byte{45, 150, 188, 218, 203, 190, 65, 56, 44, 162, 62, 82, 227, 210, 25, 108, 186, 101, 231, 161, 172, 210, 9, 223, 201, 92, 107, 50, 182, 161, 138, 147}
	sampleNickname := "bonjour"
	sampleNickname2 := "unit-test"
	sampleNickname3 := "version-0"
	nicknameRequiredFieldMissing := "required-fields-missing"
	nicknameOldLocation := "old-location-account"
	nicknameOnlyRequiredFields := "only-required-fields"
	nicknameNew := "new-account"

	createAccount := func(nickname string) *account.Account {
		acc, err := account.New(
			uint8(account.LastVersion),
			nickname,
			&types.Address{
				Object: &object.Object{
					Kind:    object.UserAddress,
					Version: types.AddressLastVersion,
					Data:    sampleAddressData,
				},
			},
			sampleSalt,
			sampleNonce,
			&types.EncryptedPrivateKey{
				Object: &object.Object{
					Kind:    object.EncryptedPrivateKey,
					Version: types.EncryptedPrivateKeyLastVersion,
					Data:    samplePrivateKeyData,
				},
			},
			&types.PublicKey{
				Object: &object.Object{
					Kind:    object.PublicKey,
					Version: types.PublicKeyLastVersion,
					Data:    samplePublicKeyData,
				},
			},
		)
		assert.NoError(t, err)

		return acc
	}

	sampleAccount := createAccount(sampleNickname)

	t.Run("Create Wallet", func(t *testing.T) {
		newWallet, err := New(walletPath)
		assert.NoError(t, err)
		w = newWallet
		assert.NotNil(t, w)
	})

	t.Run("Add Account", func(t *testing.T) {
		err := w.AddAccount(sampleAccount, true)
		assert.NoError(t, err)

		assert.Equal(t, 1, w.GetAccountCount())
		accountPath, err := w.AccountPath(sampleNickname)
		assert.NoError(t, err)
		assert.FileExists(t, accountPath)
	})

	t.Run("Add Account: nickname not unique", func(t *testing.T) {
		err := w.AddAccount(sampleAccount, true)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNicknameNotUnique)

		assert.Equal(t, 1, w.GetAccountCount())
	})

	t.Run("Add Account: address not unique", func(t *testing.T) {
		sampleAccount := createAccount(nicknameNew)
		err = w.AddAccount(sampleAccount, true)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrAddressNotUnique)
		assertAccountIsPresent(t, w, sampleNickname)
		assert.Equal(t, 1, w.GetAccountCount())

		sampleAccount.Nickname = sampleNickname
	})

	t.Run("Get Account", func(t *testing.T) {
		acc, err := w.GetAccount(sampleNickname)
		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, uint8(1), acc.Version)
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

		accountPath, err := w.AccountPath(sampleNickname2)
		assert.NoError(t, err)

		copy(t, "../../tests/wallet_unit-test.yaml", accountPath)

		assertAccountIsPresent(t, w, sampleNickname)
		acc := assertAccountIsPresent(t, w, sampleNickname2)
		assert.Equal(t, uint8(1), acc.Version)
		assert.Equal(t, 2, w.GetAccountCount())
	})

	t.Run("Invalid or unsupported version", func(t *testing.T) {
		accountPath, err := w.AccountPath(sampleNickname3)
		assert.NoError(t, err)
		copy(t, "../../tests/wallet_version-0.yaml", accountPath)
		newWallet, err := New(walletPath)
		assert.NoError(t, err)
		assertAccountIsPresent(t, w, sampleNickname)
		assertAccountIsPresent(t, newWallet, sampleNickname2)
		assert.Equal(t, 2, newWallet.GetAccountCount())
		assert.Len(t, newWallet.InvalidAccountNicknames, 1)
	})

	t.Run("Delete Account", func(t *testing.T) {
		err := w.DeleteAccount(sampleNickname)
		assert.NoError(t, err)
		assert.Equal(t, 1, w.GetAccountCount())

		accountPath, err := w.AccountPath(sampleNickname)
		assert.NoError(t, err)
		assert.NoFileExists(t, accountPath)
	})

	t.Run("Delete Account: not found", func(t *testing.T) {
		err := w.DeleteAccount("unknown")
		assert.Error(t, err)
		assert.Equal(t, 1, w.GetAccountCount())
	})

	t.Run("New wallet to discover created accounts", func(t *testing.T) {
		newWallet, err := New(walletPath)
		assert.NoError(t, err)
		assert.NotNil(t, newWallet)
		assert.Equal(t, 1, newWallet.GetAccountCount())
		assertAccountIsPresent(t, newWallet, "unit-test")
	})

	t.Run("Invalid or unsupported version: missing required fields", func(t *testing.T) {
		ClearAccounts(t, walletPath)
		accountPath, err := w.AccountPath(nicknameRequiredFieldMissing)
		assert.NoError(t, err)
		copy(t, "../../tests/wallet_required-fields-missing.yaml", accountPath)
		newWallet, err := New(walletPath)
		assert.NoError(t, err)
		assert.Equal(t, 0, newWallet.GetAccountCount())
		assert.Len(t, newWallet.InvalidAccountNicknames, 1)
	})

	t.Run("Retro-compatibility: old wallet file location", func(t *testing.T) {
		// prepare
		ClearAccounts(t, walletPath)
		accountPath, err := w.AccountPath(nicknameOldLocation)
		assert.NoError(t, err)
		copy(t, "../../tests/wallet_old-location-account.yaml", accountPath)
		// execute
		newWallet, err := New(walletPath)
		assert.NoError(t, err)

		// assert
		assert.Equal(t, 1, newWallet.GetAccountCount())
		assertAccountIsPresent(t, newWallet, nicknameOldLocation)
		assert.Len(t, newWallet.InvalidAccountNicknames, 0)
	})

	t.Run("Load account with only required fields (no address, no nickname)", func(t *testing.T) {
		// prepare
		ClearAccounts(t, walletPath)
		accountPath, err := w.AccountPath(nicknameOnlyRequiredFields)
		assert.NoError(t, err)
		copy(t, "../../tests/wallet_only-required-fields.yaml", accountPath)

		// execute
		newWallet, err := New(walletPath)
		assert.NoError(t, err)

		// assert
		assert.Equal(t, 1, newWallet.GetAccountCount())
		acc := assertAccountIsPresent(t, newWallet, nicknameOnlyRequiredFields)
		textAddress, err := acc.Address.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, addressOnlyRequiredFields, string(textAddress))
	})

	t.Run("Generate new account", func(t *testing.T) {
		acc, err := w.GenerateAccount(samplePassword, nicknameNew)
		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, 2, w.GetAccountCount())
		assertAccountIsPresent(t, w, nicknameNew)
		assertAccountIsPresent(t, w, nicknameOnlyRequiredFields)
	})

	t.Run("account from address", func(t *testing.T) {
		ClearAccounts(t, walletPath)
		err := w.Discover()
		assert.Equal(t, 0, w.GetAccountCount())
		assert.NoError(t, err)
		acc := createAccount(nicknameNew)
		err = w.AddAccount(acc, true)
		assert.NoError(t, err)
		acc, err = w.GetAccountFromAddress("not an address")
		assert.Error(t, err)
		assert.Nil(t, acc)
		acc, err = w.GetAccountFromAddress(sampleAddressText)
		assert.NoError(t, err)
		assert.NotNil(t, acc)
		assert.Equal(t, nicknameNew, acc.Nickname)
	})

	ClearAccounts(t, walletPath)
}

func copy(t *testing.T, src string, dst string) {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	assert.NoError(t, err)
	// Write data to dst
	err = os.WriteFile(dst, data, 0o644)
	assert.NoError(t, err)
}

func assertAccountIsPresent(t *testing.T, w *Wallet, nickname string) account.Account {
	acc, err := w.GetAccount(nickname)
	assert.NoError(t, err)
	assert.Equal(t, acc.Nickname, nickname)

	return *acc
}
