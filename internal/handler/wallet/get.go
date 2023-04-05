package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func HandleGet(params operations.RestWalletGetParams) middleware.Responder {
	wlt, err := wallet.Load(params.Nickname)
	if err == wallet.ErrorAccountNotFound {
		return operations.NewRestWalletGetNotFound().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	} else if err != nil {
		return operations.NewRestWalletGetBadRequest().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	}

	modelWallet := createModelWallet(*wlt)

	return operations.NewRestWalletGetOK().WithPayload(&modelWallet)
}

// HandleList handles a list request
func HandleList(params operations.RestWalletListParams) middleware.Responder {
	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewRestWalletListBadRequest().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	}

	var wlts []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		modelWallet := createModelWallet(wallets[i])
		wlts = append(wlts, &modelWallet)
	}

	return operations.NewRestWalletListOK().WithPayload(wlts)
}

func createModelWallet(wlt wallet.Wallet) models.Wallet {
	return models.Wallet{
		Nickname: wlt.Nickname,
		Address:  wlt.Address,
		KeyPair:  models.WalletKeyPair{},
	}
}
