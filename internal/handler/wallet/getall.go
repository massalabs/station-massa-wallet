package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
)

func NewGetAll(massaClient network.NodeFetcherInterface) operations.AccountListHandler {
	return &walletGetAll{massaClient: massaClient}
}

type walletGetAll struct {
	massaClient network.NodeFetcherInterface
}

func (h *walletGetAll) Handle(params operations.AccountListParams) middleware.Responder {
	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewAccountListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: err.Error(),
			})
	}
	walletsWithError, walletsWithoutError := splitWalletsPerReadError(wallets, err)
	var wlts []*models.Account

	infos, err := h.massaClient.GetAccountsInfos(walletsWithoutError)
	if err != nil {
		errMsg := "Unable to retrieve accounts infos"
		fmt.Printf("%s: %v", errMsg, err)
		return operations.NewAccountListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: errMsg,
			})
	}

	for i := 0; i < len(walletsWithoutError); i++ {
		modelWallet := createModelWallet(walletsWithoutError[i])
		modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[i].CandidateBalance))
		modelWallet.Balance = models.Amount(fmt.Sprint(infos[i].Balance))
		modelWallet.Status = wallet.StatusOK
		wlts = append(wlts, &modelWallet)
	}
	for u := 0; u < len(walletsWithError); u++ {
		status := wallet.StatusCorrupted
		modelWalletErr := createModelWallet(walletsWithError[u])
		modelWalletErr.Status = status
		wlts = append(wlts, &modelWalletErr)
	}

	return operations.NewAccountListOK().WithPayload(wlts)
}

func splitWalletsPerReadError(wallets []wallet.Wallet, err error) ([]wallet.Wallet, []wallet.Wallet) {
	var (
		walletsWithError    []wallet.Wallet
		walletsWithoutError []wallet.Wallet
	)
	for _, w := range wallets {
		if w.Status == wallet.StatusOK {
			walletsWithoutError = append(walletsWithoutError, w)
		} else {
			walletsWithError = append(walletsWithError, w)
		}
	}
	return walletsWithError, walletsWithoutError
}
