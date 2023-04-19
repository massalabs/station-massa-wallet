package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// HandleCreate handles a create request
func HandleCreate(params operations.RestWalletCreateParams) middleware.Responder {
	if len(params.Body.Nickname) == 0 {
		return operations.NewRestWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	if len(params.Body.Password) == 0 {
		return operations.NewRestWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNoPassword,
				Message: "Error: password field is mandatory.",
			})
	}

	newWallet, err := wallet.Generate(string(params.Body.Nickname), params.Body.Password)
	if err != nil {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	return New(newWallet)
}

func New(newWallet *wallet.Wallet) middleware.Responder {
	return operations.NewRestWalletCreateOK().WithPayload(
		&models.Wallet{
			Nickname: models.Nickname(newWallet.Nickname),
			Address:  newWallet.Address,
			KeyPair: models.WalletKeyPair{
				PrivateKey: "",
				PublicKey:  newWallet.GetPupKey(),
				Salt:       "",
				Nonce:      "",
			},
		})
}
