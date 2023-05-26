package wallet

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/utils"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/thyra/pkg/node/sendoperation/sellrolls"
)

func NewTradeRolls(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.TradeRollsHandler {
	return &tradeRolls{prompterApp: prompterApp, massaClient: massaClient}
}

type tradeRolls struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (t *tradeRolls) Handle(params operations.TradeRollsParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	// convert amount to uint64
	amount, err := strconv.ParseUint(string(params.Body.Amount), 10, 64)
	if err != nil {
		return operations.NewTradeRollsBadRequest().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: "Error during amount conversion",
			})
	}

	// convert fee to uint64
	fee, err := strconv.ParseUint(string(params.Body.Fee), 10, 64)
	if err != nil {
		return operations.NewTradeRollsBadRequest().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: "Error during fee conversion",
			})
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.TradeRolls,
		Msg:    fmt.Sprintf("%s %s rolls , with fee %s nonaMassa", *params.Body.Side, string(params.Body.Amount), string(params.Body.Fee)),
	}

	_, err = prompt.WakeUpPrompt(t.prompterApp, promptRequest, wlt)
	if err != nil {
		return operations.NewTradeRollsUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	operation, err := doTradeRolls(wlt, amount, fee, *params.Body.Side, t.massaClient)
	if err != nil {
		errStr := fmt.Sprintf("error %sing rolls coin: %v", *params.Body.Side, err.Error())
		t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: utils.ErrNetwork})
		return operations.NewTradeRollsInternalServerError().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: errStr,
			})
	}

	t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true, CodeMessage: utils.MsgRollTradeSuccess})
	return operations.NewTradeRollsOK().WithPayload(
		&models.OperationResponse{
			OperationID: operation.OperationID,
		})
}

func doTradeRolls(wlt *wallet.Wallet, amount, fee uint64, side string, massaClient network.NodeFetcherInterface) (*sendOperation.OperationResponse, error) {
	var operation sendOperation.Operation
	if side == "buy" {
		operation = buyrolls.New(amount)
	} else {
		operation = sellrolls.New(amount)
	}

	return network.SendOperation(wlt, massaClient, operation, fee)
}
