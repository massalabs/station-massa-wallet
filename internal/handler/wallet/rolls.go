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
	"github.com/massalabs/station/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/sellrolls"
)

func NewTradeRolls(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.TradeRollsHandler {
	return &tradeRolls{prompterApp: prompterApp, massaClient: massaClient}
}

type tradeRolls struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (t *tradeRolls) Handle(params operations.TradeRollsParams) middleware.Responder {
	acc, errResp := loadAccount(t.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	// convert amount to uint64
	amount, err := strconv.ParseUint(string(params.Body.Amount), 10, 64)
	if err != nil {
		return newErrorResponse("Error during amount conversion", errorTradeRoll, http.StatusBadRequest)
	}

	// convert fee to uint64
	fee, err := strconv.ParseUint(string(params.Body.Fee), 10, 64)
	if err != nil {
		return newErrorResponse("Error during fee conversion", errorTradeRoll, http.StatusBadRequest)
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("%s %s rolls , with fee %s nonaMassa", *params.Body.Side, string(params.Body.Amount), string(params.Body.Fee)),
	}

	promptOutput, err := prompt.WakeUpPrompt(t.prompterApp, promptRequest, acc)
	if err != nil {
		return operations.NewTradeRollsUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	output, ok := promptOutput.(*walletapp.SignPromptOutput)
	if !ok {
		return newErrorResponse(fmt.Sprintf("prompting password for roll: %v", utils.ErrInvalidInputType.Error()), utils.ErrInvalidInputType.Error(), http.StatusInternalServerError)
	}

	password := output.Password

	operation, err := doTradeRolls(acc, password, amount, fee, *params.Body.Side, t.massaClient)
	if err != nil {
		msg := fmt.Sprintf("error %sing rolls coin: %v", *params.Body.Side, err.Error())

		t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false})

		return newErrorResponse(msg, errorTransferCoin, http.StatusInternalServerError)
	}

	t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	return operations.NewTradeRollsOK().WithPayload(
		&models.OperationResponse{
			OperationID: operation.OperationID,
		})
}

func doTradeRolls(
	acc *account.Account,
	password *memguard.LockedBuffer,
	amount, fee uint64,
	side string,
	massaClient network.NodeFetcherInterface,
) (*sendOperation.OperationResponse, error) {
	var operation sendOperation.Operation
	if side == "buy" {
		operation = buyrolls.New(amount)
	} else {
		operation = sellrolls.New(amount)
	}

	return network.SendOperation(acc, password, massaClient, operation, fee)
}
