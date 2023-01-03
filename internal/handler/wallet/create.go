package wallet

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"

	"github.com/btcsuite/btcutil/base58"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
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

	newWallet, err := wallet.Generate(params.Body.Nickname, params.Body.Password)
	if err != nil {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	privK := base58.CheckEncode(newWallet.KeyPair.PrivateKey, wallet.Base58Version)
	pubK := base58.CheckEncode(newWallet.KeyPair.PublicKey, wallet.Base58Version)
	salt := base58.CheckEncode(newWallet.KeyPair.Salt[:], wallet.Base58Version)
	nonce := base58.CheckEncode(newWallet.KeyPair.Nonce[:], wallet.Base58Version)

	return operations.NewRestWalletCreateOK().WithPayload(
		&models.Wallet{
			Nickname: newWallet.Nickname,
			Address:  newWallet.Address,
			KeyPair: models.WalletKeyPair{
				PrivateKey: privK,
				PublicKey:  pubK,
				Salt:       salt,
				Nonce:      nonce,
			},
		})
}
