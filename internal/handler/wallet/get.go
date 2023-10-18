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

	modelWallet, err := newAccountModel(*acc)
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
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

		password, _ := promptOutput.(*memguard.LockedBuffer)

		guardedPrivateKey, err := acc.PrivateKeyTextInClear(password)
		if err != nil {
			return operations.NewGetAccountInternalServerError().WithPayload(&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
		}

		defer guardedPrivateKey.Destroy()

		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: true})

		modelWallet.KeyPair, err = newKeyPairModel(*acc, guardedPrivateKey)
		if err != nil {
			return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
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

func newAccountModel(acc account.Account) (*models.Account, error) {
	address, err := acc.Address.MarshalText()
	if err != nil {
		return nil, err
	}

	return &models.Account{
		Nickname: models.Nickname(acc.Nickname),
		Address:  models.Address(address),
		KeyPair:  models.KeyPair{},
		Status:   accountStatusOK,
	}, nil
}

func newKeyPairModel(acc account.Account, guardedPrivateKey *memguard.LockedBuffer) (models.KeyPair, error) {
	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return models.KeyPair{}, err
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
