package wallet

import (
	"fmt"
	"net/http"

	"github.com/awnumar/memguard"
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

const (
	accountStatusOK        = "ok"
	accountStatusCorrupted = "corrupted"
)

// NewGet instantiates a Get Handler
// The "classical" way is not possible because we need to pass to the handler a wallet.WalletPrompterInterface.
func NewGet(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.GetAccountHandler {
	return &walletGet{prompterApp: prompterApp, massaClient: massaClient}
}

type walletGet struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (w *walletGet) Handle(params operations.GetAccountParams) middleware.Responder {
	acc, resp := loadAccount(w.prompterApp.App().WalletManager, params.Nickname)
	if resp != nil || acc == nil {
		return resp
	}

	modelWallet, resp := newAccountModel(*acc)
	if resp != nil {
		return resp
	}

	// if request not ciphered data, ask for password and unprotect the wallet
	if params.Ciphered != nil && !*params.Ciphered {
		promptRequest := prompt.PromptRequest{
			Action: walletapp.Unprotect,
			Msg:    fmt.Sprintf("Unprotect wallet %s", acc.Nickname),
		}

		promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, promptRequest, acc)
		if err != nil {
			return operations.NewGetAccountUnauthorized().WithPayload(
				&models.Error{
					Code:    fmt.Sprint(http.StatusUnauthorized),
					Message: "Unable to unprotect wallet",
				})
		}

		guardedPassword, _ := promptOutput.(*memguard.LockedBuffer)
		guardedPrivateKey, err := acc.PrivateKeyTextInClear(guardedPassword)
		if err != nil {
			return operations.NewGetAccountInternalServerError().WithPayload(&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
		}

		defer guardedPrivateKey.Destroy()

		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountUnprotected})

		modelWallet.KeyPair, resp = newKeyPairModel(*acc, guardedPrivateKey)
		if resp != nil {
			return resp
		}
	}

	infos, err := w.massaClient.GetAccountsInfos([]account.Account{*acc})
	if err != nil {
		return operations.NewGetAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    utils.ErrNetwork,
				Message: "Unable to retrieve account infos",
			})
	}

	modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[0].CandidateBalance))
	modelWallet.Balance = models.Amount(fmt.Sprint(infos[0].Balance))

	return operations.NewGetAccountOK().WithPayload(modelWallet)
}

func newAccountModel(acc account.Account) (*models.Account, middleware.Responder) {
	address, err := acc.Address.MarshalText()
	if err != nil {
		return nil, operations.NewGetAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallet,
				Message: ErrorAddressInvalid.Error(),
			})
	}

	return &models.Account{
		Nickname: models.Nickname(acc.Nickname),
		Address:  models.Address(address),
		KeyPair:  models.KeyPair{},
		Status:   accountStatusOK,
	}, nil
}

func newKeyPairModel(acc account.Account, guardedPrivateKey *memguard.LockedBuffer) (models.KeyPair, middleware.Responder) {
	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return models.KeyPair{}, operations.NewGetAccountInternalServerError().WithPayload(&models.Error{
			Code:    errorGetWallets,
			Message: err.Error(),
		})
	}

	return models.KeyPair{
		PrivateKey: string(guardedPrivateKey.Bytes()),
		PublicKey:  string(publicKeyBytes),
		Salt:       stringifyByteArray(acc.Salt[:]),
		Nonce:      stringifyByteArray(acc.Nonce[:]),
	}, nil
}

func stringifyByteArray(b []byte) string {
	return base58.CheckEncode(b, 0x00)
}
