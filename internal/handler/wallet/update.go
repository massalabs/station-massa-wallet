package wallet

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

type PromptRequestUpdateAccountData struct {
	Nickname string
	Balance  string
}

func NewUpdateAccount(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.UpdateAccountHandler {
	return &walletUpdateAccount{prompterApp: prompterApp, massaClient: massaClient}
}

type walletUpdateAccount struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

// HandleDelete handles an update request
func (w *walletUpdateAccount) Handle(params operations.UpdateAccountParams) middleware.Responder {
	acc, resp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if resp != nil {
		return resp
	}

	newAcc, err := w.handleUpdateAccount(acc, params.Body.Nickname)
	if err != nil {
		return newErrorResponse(err.Error(), "", http.StatusInternalServerError)
	}

	modelWallet, err := newAccountModel(*newAcc)
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	infos, err := w.massaClient.GetAccountsInfos([]account.Account{*acc})
	if err != nil {
		return operations.NewGetAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve account infos",
			})
	}

	modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[0].CandidateBalance))
	modelWallet.Balance = models.Amount(fmt.Sprint(infos[0].Balance))

	return operations.NewGetAccountOK().WithPayload(modelWallet)
}

func (w *walletUpdateAccount) handleUpdateAccount(acc *account.Account, newNickname models.Nickname) (*account.Account, error) {
	// check if the nickname does not change
	if acc.Nickname == string(newNickname) {
		return nil, wallet.ErrNicknameNotUnique
	}

	// save the old nickname in a variable
	oldNickname := acc.Nickname

	newAcc, err := account.New(
		acc.Version,
		string(newNickname),
		acc.Address,
		acc.Salt,
		acc.Nonce,
		acc.CipheredData,
		acc.PublicKey,
	)
	if err != nil {
		return nil, fmt.Errorf("creating new account: %w", err)
	}

	err = w.prompterApp.App().Wallet.DeleteAccount(string(oldNickname))
	if err != nil {
		return nil, err
	}

	err = w.prompterApp.App().Wallet.AddAccount(newAcc, true)
	if err != nil {
		return nil, err
	}

	return newAcc, nil
}
