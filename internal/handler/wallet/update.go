package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
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
func (m *walletUpdateAccount) Handle(params operations.UpdateAccountParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	newWlt, errModify := m.handleUpdateAccount(wlt, params.Body.Nickname)
	if errModify != nil {
		return operations.NewGetAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errModify.CodeErr,
				Message: errModify.Err.Error(),
			})
	}

	modelWallet := createModelWallet(*newWlt)

	infos, err := m.massaClient.GetAccountsInfos([]wallet.Wallet{*wlt})
	if err != nil {
		return operations.NewGetAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve account infos",
			})
	}

	modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[0].CandidateBalance))
	modelWallet.Balance = models.Amount(fmt.Sprint(infos[0].Balance))

	return operations.NewGetAccountOK().WithPayload(&modelWallet)
}

func (m *walletUpdateAccount) handleUpdateAccount(wlt *wallet.Wallet, newNickname models.Nickname) (*wallet.Wallet, *wallet.WalletError) {
	// check if the nickname does not change
	if wlt.Nickname == string(newNickname) {
		return nil, &wallet.WalletError{Err: fmt.Errorf("nickname is the same"), CodeErr: utils.ErrSameNickname}
	}

	// Validate nickname
	if !wallet.NicknameIsValid(string(newNickname)) {
		return nil, &wallet.WalletError{Err: fmt.Errorf("invalid nickname"), CodeErr: utils.ErrInvalidNickname}
	}

	// Validate nickname uniqueness
	err := wallet.NicknameIsUnique(string(newNickname))
	if err != nil {
		return nil, &wallet.WalletError{Err: err, CodeErr: utils.ErrDuplicateNickname}
	}

	oldNickname := wlt.Nickname

	// persist new nickname before deleting old file
	wlt.Nickname = string(newNickname)
	err = wlt.Persist()
	if err != nil {
		return nil, &wallet.WalletError{
			Err:     fmt.Errorf("persisting the modified account: %w", err),
			CodeErr: utils.ErrAccountFile,
		}
	}

	// delete old file
	if wallet.DeleteAccount(oldNickname) != nil {
		return nil, &wallet.WalletError{
			Err:     fmt.Errorf("persisting the old account: %w", err),
			CodeErr: utils.ErrAccountFile,
		}
	}

	return wlt, nil
}
