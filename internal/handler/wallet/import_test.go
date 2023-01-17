package wallet

import (
	"testing"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/guiModal"
)

func Test_walletImport_Handle(t *testing.T) {

	var mockWalletInfoModal guiModal.WalletInfoAsker
	wimport := NewImport(mockWalletInfoModal)
	params := operations.RestWalletImportParams{}

	response := wimport.Handle(params)
	if response == ImportWalletError(errorImportWalletCanceled, errorImportWalletCanceled) {
		t.Error("Wallet Canceled", response)
	}

	if response == ImportWalletError("error", "Error: a wallet with the same nickname already exists.") {
		t.Error("Name Already Taken", response)
	}

}
