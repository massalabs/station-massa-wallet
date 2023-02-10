package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// HandleList handles a list request
func HandleList(params operations.RestWalletListParams) middleware.Responder {
	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewRestWalletListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	}

	var wlts []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		wlts = append(wlts,
			&models.Wallet{
				Nickname: wallets[i].Nickname,
				Address:  wallets[i].Address,
				KeyPair:  models.WalletKeyPair{},
			})
	}

	return operations.NewRestWalletListOK().WithPayload(wlts)
}
