package wallet

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/awnumar/memguard"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
)

type PromptRequestTransferData struct {
	NicknameFrom     string
	Amount           string
	Fee              string
	RecipientAddress string
}

func NewTransferCoin(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.TransferCoinHandler {
	return &transferCoin{prompterApp: prompterApp, massaClient: massaClient}
}

type transferCoin struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (t *transferCoin) Handle(params operations.TransferCoinParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	acc, resp := loadAccount(t.prompterApp.App().WalletManager, params.Nickname)
	if resp != nil {
		return resp
	}

	// convert amount to uint64
	amount, err := strconv.ParseUint(string(params.Body.Amount), 10, 64)
	if err != nil {
		return operations.NewTransferCoinBadRequest().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: "Error during amount conversion",
			})
	}

	// convert fee to uint64
	fee, err := strconv.ParseUint(string(params.Body.Fee), 10, 64)
	if err != nil {
		return operations.NewTransferCoinBadRequest().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: "Error during fee conversion",
			})
	}

	promptRequest := prompt.PromptRequest{
		Action:      walletapp.Transfer,
		CodeMessage: utils.MsgTransferRequest,
		Data: PromptRequestTransferData{
			NicknameFrom:     acc.Nickname,
			Amount:           string(params.Body.Amount),
			Fee:              string(params.Body.Fee),
			RecipientAddress: *params.Body.RecipientAddress,
		},
	}

	promptOutput, err := prompt.WakeUpPrompt(t.prompterApp, promptRequest, acc)
	if err != nil {
		return operations.NewTransferCoinUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	guardedPassword, _ := promptOutput.(*memguard.LockedBuffer)

	// create the transaction and send it to the network
	operation, err := doTransfer(acc, guardedPassword, amount, fee, *params.Body.RecipientAddress, t.massaClient)
	if err != nil {
		msg := fmt.Sprintf("error transferring coin: %v", err.Error())

		t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: errorTransferCoin})

		return newErrorResponse(msg, errorTransferCoin, http.StatusInternalServerError)
	}

	t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgTransferSuccess})

	return operations.NewTransferCoinOK().WithPayload(
		&models.OperationResponse{
			OperationID: operation.OperationID,
		})
}

func doTransfer(
	acc *account.Account,
	guardedPassword *memguard.LockedBuffer,
	amount, fee uint64,
	recipientAddress string,
	massaClient network.NodeFetcherInterface,
) (*sendOperation.OperationResponse, error) {
	operation, err := transaction.New(recipientAddress, amount)
	if err != nil {
		return nil, fmt.Errorf("Error during transaction creation: %w", err)
	}

	return network.SendOperation(acc, guardedPassword, massaClient, operation, fee)
}
