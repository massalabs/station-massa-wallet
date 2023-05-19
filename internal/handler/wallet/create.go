package wallet

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func NewCreateAccount(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.CreateAccountHandler {
	return &walletCreate{prompterApp: prompterApp, massaClient: massaClient}
}

type walletCreate struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (w *walletCreate) Handle(params operations.CreateAccountParams) middleware.Responder {
	nickname := strings.TrimSpace(string(params.Nickname))

	if len(nickname) == 0 {
		return operations.NewCreateAccountBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNoNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.NewPassword,
		Msg:    "Define a password",
		Data:   nil,
	}

	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, promptRequest, nil)
	if err != nil {
		return operations.NewCreateAccountUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to create wallet",
			})
	}

	password, _ := promptOutput.(*string)

	wlt, err := wallet.Generate(nickname, *password)
	if err != nil {
		errStr := fmt.Sprintf("Unable to create account: %v", err)
		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, Data: errStr})

		return operations.NewCreateAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, Data: "New password created"})

	infos, err := w.massaClient.GetAccountsInfos([]wallet.Wallet{*wlt})
	if err != nil {
		return operations.NewCreateAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve accounts infos",
			})
	}

	return operations.NewCreateAccountOK().WithPayload(
		&models.Account{
			Nickname:         models.Nickname(wlt.Nickname),
			Address:          wlt.Address,
			CandidateBalance: models.Amount(fmt.Sprint(infos[0].CandidateBalance)),
			Balance:          models.Amount(fmt.Sprint(infos[0].Balance)),
			KeyPair: models.AccountKeyPair{
				PrivateKey: "",
				PublicKey:  wlt.GetPupKey(),
				Salt:       "",
				Nonce:      "",
			},
		})
}
