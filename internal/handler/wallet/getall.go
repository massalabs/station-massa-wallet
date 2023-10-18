package wallet

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station/pkg/logger"
)

func NewGetAll(wallet *wallet.Wallet, massaClient network.NodeFetcherInterface) operations.AccountListHandler {
	return &walletGetAll{wallet: wallet, massaClient: massaClient}
}

type walletGetAll struct {
	wallet      *wallet.Wallet
	massaClient network.NodeFetcherInterface
}

func (w *walletGetAll) Handle(params operations.AccountListParams) middleware.Responder {
	err := w.wallet.Discover()
	if err != nil {
		errMsg := "Unable to discover accounts"
		logger.Errorf("%s: %v", errMsg, err)

		return operations.NewAccountListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: errMsg,
			})
	}

	accounts := w.wallet.AllAccounts()

	infos, err := w.massaClient.GetAccountsInfos(accounts)
	if err != nil {
		errMsg := "Unable to retrieve accounts infos"
		logger.Warnf("%s: %v", errMsg, err)

		return operations.NewAccountListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: errMsg,
			})
	}

	var accountModels []*models.Account

	for i := 0; i < len(accounts); i++ {
		modelWallet, err := newAccountModel(accounts[i])
		if err != nil {
			return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
		}
		modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[i].CandidateBalance))
		modelWallet.Balance = models.Amount(fmt.Sprint(infos[i].Balance))
		accountModels = append(accountModels, modelWallet)
	}

	for u := 0; u < len(w.wallet.InvalidAccountNicknames); u++ {
		invalidAccount := &models.Account{
			Nickname: models.Nickname(w.wallet.InvalidAccountNicknames[u]),
		}
		invalidAccount.Status = accountStatusCorrupted
		accountModels = append(accountModels, invalidAccount)
	}

	return operations.NewAccountListOK().WithPayload(accountModels)
}
