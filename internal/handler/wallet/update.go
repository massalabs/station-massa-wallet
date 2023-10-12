package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station-massa-wallet/pkg/walletmanager"
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
	acc, resp := loadAccount(w.prompterApp.App().WalletManager, params.Nickname)
	if resp != nil {
		return resp
	}

	newAcc, errModify := w.handleUpdateAccount(acc, params.Body.Nickname)
	if errModify != nil {
		return operations.NewGetAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errModify.CodeErr,
				Message: errModify.Err.Error(),
			})
	}

	modelWallet, resp := newAccountModel(*newAcc)
	if resp != nil {
		return resp
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

func (w *walletUpdateAccount) handleUpdateAccount(acc *account.Account, newNickname models.Nickname) (*account.Account, *walletmanager.WalletError) {
	// check if the nickname does not change
	if acc.Nickname == string(newNickname) {
		return nil, &walletmanager.WalletError{Err: fmt.Errorf("nickname is the same"), CodeErr: utils.ErrSameNickname}
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
		return nil, &walletmanager.WalletError{Err: err, CodeErr: utils.ErrInvalidNickname}
	}

	err = w.prompterApp.App().WalletManager.DeleteAccount(string(oldNickname))
	if err != nil {
		return nil, &walletmanager.WalletError{Err: err, CodeErr: utils.WailsErrorCode(err)}
	}

	err = w.prompterApp.App().WalletManager.AddAccount(newAcc, true)
	if err != nil {
		return nil, &walletmanager.WalletError{Err: err, CodeErr: utils.WailsErrorCode(err)}
	}

	return newAcc, nil
}
