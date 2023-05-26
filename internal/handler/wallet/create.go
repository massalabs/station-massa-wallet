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
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
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

	wlt, errGenerate := wallet.Generate(nickname, *password)
	if errGenerate != nil {
		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: errGenerate.CodeErr})

		// At this stage, we can't know if its 400 or 500 (let's say 400 because in the test case 400 make sense)
		return operations.NewCreateAccountBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: errGenerate.Err.Error(),
			})
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountCreated})

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
