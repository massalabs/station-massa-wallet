package wallet

import (
	"testing"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

type MockWalletInfoModal struct{}

func (m *MockWalletInfoModal) WalletInfo() (string, string, string, error) {
	return "password", "walletName", "privateKey", nil
}

func Test_walletImport_Handle(t *testing.T) {

	mockWalletInfoModal := &MockWalletInfoModal{}

	importHandler := NewImport(mockWalletInfoModal)

	// Import wallet canceled
	resp := importHandler.Handle(operations.NewRestWalletImportParams())
	if resp == nil || resp != ImportWalletError(errorImportWalletCanceled, errorImportWalletCanceled) {
		t.Error("Expected error 'Error: wallet import canceled', got", resp)
	}

	// NickName already taken
	_, err := wallet.Load("nickNameAlreadyTaken")
	if err == nil {
		t.Error("Expected error 'Error: wallet import canceled', got", resp)
	}

	// Wallet well imported
	_, err = wallet.Imported("validNickName", "validPK", "validPassword")
	if err == nil {
		t.Error("UnexpectedError, got", err)
	}

}
