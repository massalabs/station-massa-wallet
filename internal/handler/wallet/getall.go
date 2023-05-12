package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
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

	var wlts []*models.Account

	infos, err := h.massaClient.GetAccountsInfos(wallets)
	if err != nil {
		return operations.NewAccountListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve accounts infos",
			})
	}

	for i := 0; i < len(wallets); i++ {
		modelWallet := createModelWallet(wallets[i])
		modelWallet.CandidateBalance = models.Amount(fmt.Sprint(infos[i].CandidateBalance))
		modelWallet.Balance = models.Amount(fmt.Sprint(infos[i].Balance))
		wlts = append(wlts, &modelWallet)
	}

	return operations.NewAccountListOK().WithPayload(wlts)
}
