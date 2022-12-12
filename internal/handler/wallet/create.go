package wallet

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/base58"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

//nolint:nolintlint,ireturn
func NewCreate(walletStorage *sync.Map) operations.RestWalletCreateHandler {
	return &walletCreate{walletStorage: walletStorage}
}

type walletCreate struct {
	walletStorage *sync.Map
}

//nolint:nolintlint,ireturn,funlen
func (c *walletCreate) Handle(params operations.RestWalletCreateParams) middleware.Responder {
	if params.Body.Nickname == nil || len(*params.Body.Nickname) == 0 {
		return operations.NewRestWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	_, ok := c.walletStorage.Load(*params.Body.Nickname)
	if ok {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorAlreadyExists,
				Message: "Error: a wallet with the same nickname already exists.",
			})
	}

	if params.Body.Password == nil || len(*params.Body.Password) == 0 {
		return operations.NewRestWalletCreateBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNoPassword,
				Message: "Error: password field is mandatory.",
			})
	}

	newWallet, err := wallet.New(*params.Body.Nickname)
	if err != nil {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	return CreateNewWallet(params.Body.Nickname, params.Body.Password, c.walletStorage, newWallet)
}

//nolint:lll,nolintlint,ireturn
func CreateNewWallet(nickname *string, password *string, storage *sync.Map, newWallet *wallet.Wallet) middleware.Responder {
	err := newWallet.Protect(*password, 0)
	if err != nil {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	bytesOutput, err := json.Marshal(newWallet)
	if err != nil {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	err = os.WriteFile(wallet.GetWalletFile(*nickname), bytesOutput, fileModeUserRW)
	if err != nil {
		return operations.NewRestWalletCreateInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	storage.Store(newWallet.Nickname, newWallet)

	privK := base58.CheckEncode(newWallet.KeyPairs[0].PrivateKey)
	pubK := base58.CheckEncode(newWallet.KeyPairs[0].PublicKey)
	salt := base58.CheckEncode(newWallet.KeyPairs[0].Salt[:])
	nonce := base58.CheckEncode(newWallet.KeyPairs[0].Nonce[:])

	return operations.NewRestWalletCreateOK().WithPayload(
		&models.Wallet{
			Nickname: &newWallet.Nickname,
			Address:  &newWallet.Address,
			KeyPairs: []*models.WalletKeyPairsItems0{{
				PrivateKey: &privK,
				PublicKey:  &pubK,
				Salt:       &salt,
				Nonce:      &nonce,
			}},
			Balance: 0,
		})
}
