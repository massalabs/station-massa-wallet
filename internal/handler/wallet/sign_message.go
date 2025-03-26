package wallet

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

const (
	Message = int(-1)
)

type PromptRequestSignMessageData struct {
	Description   string
	OperationType int
	WalletAddress string
	Nickname      string
	PlainText     string
	DisplayData   bool
}

func NewSignMessage(prompterApp prompt.WalletPrompterInterface) operations.SignMessageHandler {
	return &walletSignMessage{prompterApp: prompterApp}
}

type walletSignMessage struct {
	prompterApp prompt.WalletPrompterInterface
}

func (w *walletSignMessage) Handle(params operations.SignMessageParams) middleware.Responder {
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	// Create a promptRequest for signing the message
	promptRequest, err := prepareSignMessagePromptRequest(*acc, params.Body)
	if err != nil {
		return newErrorResponse(err.Error(), errorSignDecodeMessage, http.StatusBadRequest)
	}

	output, err := PromptForOperation(w.prompterApp, acc, promptRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to unprotect wallet: %s", err.Error())
		if errors.Is(err, utils.ErrWrongPassword) || errors.Is(err, utils.ErrActionCanceled) {
			return newErrorResponse(msg, errorGetWallets, http.StatusUnauthorized)
		}

		return newErrorResponse(msg, errorGetWallets, http.StatusInternalServerError)
	}

	signature, err := acc.Sign(output.Password, []byte(params.Body.Message))
	if err != nil {
		return newErrorResponse(fmt.Sprintf("unable to sign message: %s", err.Error()), errorGetWallets, http.StatusInternalServerError)
	}

	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	cfg := config.Get()

	if cfg.HasEnabledRule(acc.Nickname) {
		err = cache.CachePrivateKeyFromPassword(acc, output.Password)
		if err != nil {
			return newErrorResponse(err.Error(), errorCachePrivateKey, http.StatusInternalServerError)
		}
	}

	// Return the signature and public key as the response
	return operations.NewSignMessageOK().WithPayload(
		&models.SignResponse{
			PublicKey: string(publicKeyBytes),
			Signature: signature,
		})
}

func prepareSignMessagePromptRequest(acc account.Account, body *models.SignMessageRequest) (*prompt.PromptRequest, error) {
	DisplayData := true
	// Check if DisplayData is provided in the request, if not, use the default (true)
	if body.DisplayData != nil {
		DisplayData = *body.DisplayData
	}

	addressBytes, err := acc.Address.MarshalText()
	if err != nil {
		return nil, err
	}

	return &prompt.PromptRequest{
		Action: walletapp.Sign,
		Data: PromptRequestSignMessageData{
			Description:   body.Description,
			OperationType: Message,
			DisplayData:   DisplayData,
			PlainText:     body.Message,
			WalletAddress: string(addressBytes),
			Nickname:      acc.Nickname,
		},
	}, nil
}
