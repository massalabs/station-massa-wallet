package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
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
	oldNickname := wlt.Nickname

	// persist new nickname before deleting old file
	wlt.Nickname = string(newNickname)
	err := wlt.Persist()
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
