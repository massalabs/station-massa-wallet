package wallet

import (
	"github.com/btcsuite/btcutil/base58"
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
		privK := base58.CheckEncode(wallets[i].KeyPair.PrivateKey, wallet.Base58Version)
		pubK := base58.CheckEncode(wallets[i].KeyPair.PublicKey, wallet.Base58Version)
		salt := base58.CheckEncode(wallets[i].KeyPair.Salt[:], wallet.Base58Version)
		nonce := base58.CheckEncode(wallets[i].KeyPair.Nonce[:], wallet.Base58Version)
		wlts = append(wlts,
			&models.Wallet{
				Nickname: wallets[i].Nickname,
				Address:  wallets[i].Address,
				KeyPair:  models.WalletKeyPair{
					PrivateKey: privK,
					PublicKey:  pubK,
					Salt:       salt,
					Nonce:      nonce,
				},
			})
	}

	return operations.NewRestWalletListOK().WithPayload(wlts)
}
