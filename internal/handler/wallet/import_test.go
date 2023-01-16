package wallet

import (
	"sync"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

func Test_walletImport_Handle(t *testing.T) {

	app := test.NewApp()
	defer test.NewApp()

	walletStorage := new(sync.Map)
	importHandler := NewImport(walletStorage, &app)

	// Import wallet canceled
	resp := importHandler.Handle(operations.NewRestWalletImportParams())
	if resp == nil || resp != NewWalletError(errorImportWalletCanceled, errorImportWalletCanceled) {
		t.Error("Expected error 'Error: wallet import canceled', got", resp)
	}

	// NickName already taken
	_, inStore := walletStorage.Load("nickNameAlreadyTaken")
	if inStore {
		t.Error("Expected error 'Error: wallet import canceled', got", resp)
	}

	// Wallet well imported
	_, err := wallet.Imported("validNickName", "validPK", "validPassword")
	if err == nil {
		t.Error("UnexpectedError, got", err)
	}

}
