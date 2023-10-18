package wallet

import (
	"fmt"
	"net/http"

	"github.com/awnumar/memguard"
	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

type PromptRequestSignMessageData struct {
	Description   string
	OperationType string
	WalletAddress string
	PlainText     string
	DisplayData   bool
}

func NewSignMessage(prompterApp prompt.WalletPrompterInterface, gc gcache.Cache) operations.SignMessageHandler {
	return &walletSignMessage{gc: gc, prompterApp: prompterApp}
}

type walletSignMessage struct {
	prompterApp prompt.WalletPrompterInterface
	gc          gcache.Cache
}

func (w *walletSignMessage) Handle(params operations.SignMessageParams) middleware.Responder {
	acc, resp := loadAccount(w.prompterApp.App().WalletManager, params.Nickname)
	if resp != nil {
		return resp
	}

	// Create a promptRequest for signing the message
	promptRequest, resp := prepareSignMessagePromptRequest(*acc, params.Body)
	if resp != nil {
		return resp
	}

	// Use the prompt-based logic to sign the message
	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, *promptRequest, acc)
	if err != nil {
		return operations.NewSignMessageUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	guardedPassword, _ := promptOutput.(*memguard.LockedBuffer)

	signature, err := acc.Sign(guardedPassword, []byte(params.Body.Message))
	if err != nil {
		return newErrorResponse(fmt.Sprintf("unable to sign message: %s", err.Error()), errorGetWallets, http.StatusInternalServerError)
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountUnprotected})

	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return newErrorResponse(fmt.Sprintf("unable to marshal public key: %s", err.Error()), errorGetWallets, http.StatusInternalServerError)
	}

	// Return the signature and public key as the response
	return operations.NewSignMessageOK().WithPayload(
		&models.SignResponse{
			PublicKey: string(publicKeyBytes),
			Signature: signature,
		})
}

func prepareSignMessagePromptRequest(acc account.Account, body *models.SignMessageRequest) (*prompt.PromptRequest, middleware.Responder) {
	DisplayData := true
	// Check if DisplayData is provided in the request, if not, use the default (true)
	if body.DisplayData != nil {
		DisplayData = *body.DisplayData
	}

	address, err := acc.Address.MarshalText()
	if err != nil {
		return nil, newErrorResponse(fmt.Sprintf("Unable to marshal address: %s", err.Error()), errorGetWallets, http.StatusInternalServerError)
	}

	return &prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", acc.Nickname),
		Data: PromptRequestSignMessageData{
			Description:   body.Description,
			OperationType: Message,
			DisplayData:   DisplayData,
			PlainText:     body.Message,
			WalletAddress: string(address),
		},
	}, nil
}
