package wallet

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func NewDelete(prompterApp prompt.WalletPrompterInterface) operations.RestAccountDeleteHandler {
	return &walletDelete{prompterApp: prompterApp}
}

type walletDelete struct {
	prompterApp prompt.WalletPrompterInterface
}

// HandleDelete handles a delete request
func (w *walletDelete) Handle(params operations.RestAccountDeleteParams) middleware.Responder {
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return operations.NewRestAccountDeleteBadRequest().WithPayload(
			&models.Error{
				Code:    errorGetWallet,
				Message: "Error cannot load wallet: " + err.Error(),
			})
	}

	go handleDelete(wallet, w.prompterApp)

	return operations.NewRestAccountDeleteNoContent()
}

func handleDelete(wlt *wallet.Wallet, prompterApp prompt.WalletPrompterInterface) {
	promptData := &prompt.PromptRequestData{
		Msg:  fmt.Sprintf("Deleting wallet %s:", wlt.Nickname),
		Data: nil,
	}

	_, err := prompt.PromptPassword(prompterApp, wlt, walletapp.Password, promptData)
	if err != nil {
		return
	}

	if wlt.DeleteFile() != nil {
		errStr := "error deleting wallet:" + err.Error()
		fmt.Println(errStr)
		prompterApp.EmitEvent(walletapp.PasswordResultEvent,
			walletapp.EventData{Success: false, Data: errStr})
	}

	prompterApp.EmitEvent(walletapp.PasswordResultEvent,
		walletapp.EventData{Success: true, Data: "Delete Success"})
}
