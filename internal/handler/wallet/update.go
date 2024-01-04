package wallet

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
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
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	newAcc, err := w.handleUpdateAccount(acc, params.Body.Nickname)
	if err != nil {
		if errors.Is(err, wallet.ErrNicknameNotUnique) {
			return newErrorResponse(err.Error(), utils.ErrDuplicateNickname, http.StatusBadRequest)
		} else {
			return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
		}
	}

	modelWallet, err := newAccountModel(newAcc)
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	infos, err := w.massaClient.GetAccountsInfos([]*account.Account{acc})
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
	nickname := string(newNickname)
	if acc.Nickname == nickname {
		return nil, wallet.ErrNicknameNotUnique
	}

	// save the old nickname in a variable
	oldNickname := acc.Nickname

	newAcc, err := account.New(
		*acc.Version,
		nickname,
		acc.Address,
		acc.Salt,
		acc.Nonce,
		acc.CipheredData,
		acc.PublicKey,
	)
	if err != nil {
		return nil, fmt.Errorf("creating new account: %w", err)
	}

	// check if the new nickname is unique
	err = w.prompterApp.App().Wallet.NicknameIsUnique(nickname)
	if err != nil {
		return nil, fmt.Errorf("checking nickname uniqueness: %w", err)
	}

	// force the account creation to bypass the address uniqueness check
	err = w.prompterApp.App().Wallet.AddAccount(newAcc, true, true)
	if err != nil {
		return nil, fmt.Errorf("adding account: %w", err)
	}

	// delete the old account after the new one has been created to avoid losing the account
	err = w.prompterApp.App().Wallet.DeleteAccount(string(oldNickname))
	if err != nil {
		return nil, fmt.Errorf("deleting account: %w", err)
	}

	return newAcc, nil
}
