package wallet

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/awnumar/memguard"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/ico"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
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
		return newErrorResponse("Unable to create wallet", errorCanceledAction, http.StatusUnauthorized)
	}

	password, _ := promptOutput.(*memguard.LockedBuffer)

	acc, err := w.prompterApp.App().Wallet.GenerateAccount(password, nickname)
	if err != nil {
		w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.WailsErrorCode(err)})

		// At this stage, we can't know if its 400 or 500 (let's say 400 because in the test case 400 make sense)
		return operations.NewCreateAccountBadRequest().WithPayload(
			&models.Error{
				Code:    errorCreateNew,
				Message: err.Error(),
			})
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	infos, err := w.massaClient.GetAccountsInfos([]*account.Account{acc})
	if err != nil {
		return operations.NewCreateAccountInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallets,
				Message: "Unable to retrieve accounts infos",
			})
	}

	address, err := acc.Address.MarshalText()
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	publicKey, err := acc.PublicKey.MarshalText()
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	//ICOQUEST: To be removed when ICO is over
	//nolint:errcheck
	ico.ValidateQuest("CREATE_WALLET", string(address))

	return operations.NewCreateAccountOK().WithPayload(
		&models.Account{
			Nickname:         models.Nickname(acc.Nickname),
			Address:          models.Address(address),
			CandidateBalance: models.Amount(fmt.Sprint(infos[0].CandidateBalance)),
			Balance:          models.Amount(fmt.Sprint(infos[0].Balance)),
			KeyPair: models.KeyPair{
				PrivateKey: "",
				PublicKey:  string(publicKey),
				Salt:       "",
				Nonce:      "",
			},
		})
}
