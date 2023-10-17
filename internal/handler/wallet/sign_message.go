package wallet

import (
	"fmt"

	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
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

func (s *walletSignMessage) Handle(params operations.SignMessageParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	// Create a promptRequest for signing the message
	promptRequest := s.prepareSignMessagePromptRequest(wlt, params.Body)

	// Use the prompt-based logic to sign the message
	_, err := prompt.WakeUpPrompt(s.prompterApp, promptRequest, wlt)
	if err != nil {
		return operations.NewSignUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	s.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountUnprotected})

	// Sign the message using the wallet
	signature := wlt.Sign(false, []byte(params.Body.Message))

	// Return the signature and public key as the response
	return operations.NewSignMessageOK().WithPayload(
		&models.SignResponse{
			PublicKey: wlt.GetPupKey(),
			Signature: signature,
		})
}

func (s *walletSignMessage) prepareSignMessagePromptRequest(wlt *wallet.Wallet, body *models.SignMessageRequest) prompt.PromptRequest {
	DisplayData := true
	// Check if DisplayData is provided in the request, if not, use the default (true)
	if body.DisplayData != nil {
		DisplayData = *body.DisplayData
	}

	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
		Data: PromptRequestSignMessageData{
			Description:   body.Description,
			OperationType: Message,
			DisplayData:   DisplayData,
			PlainText:     body.Message,
			WalletAddress: wlt.Address,
		},
	}
}
