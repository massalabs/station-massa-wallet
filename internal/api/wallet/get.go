package wallet

import (
	"sync"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-massa-core/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-core/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-massa-core/pkg/wallet"
)

//nolint:nolintlint,ireturn
func NewList(walletStorage *sync.Map) operations.RestWalletListHandler {
	return &walletList{walletStorage: walletStorage}
}

type walletList struct {
	walletStorage *sync.Map
}

//nolint:nolintlint,ireturn
func (c *walletList) Handle(params operations.RestWalletListParams) middleware.Responder {
	wallets, err := wallet.LoadAll()
	if err != nil {
		return operations.NewRestWalletListInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletGetWallets,
				Message: err.Error(),
			})
	}

	var wal []*models.Wallet

	for i := 0; i < len(wallets); i++ {
		walletss := &models.Wallet{
			Nickname: &wallets[i].Nickname,
			Address:  &wallets[i].Address,
			KeyPairs: []*models.WalletKeyPairsItems0{},
		}

		wal = append(wal, walletss)
	}

	return operations.NewRestWalletListOK().WithPayload(wal)
}
