package wallet

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/network"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/transaction"
)

func NewTransferCoin(prompterApp prompt.WalletPrompterInterface, massaClient network.NodeFetcherInterface) operations.TransferCoinHandler {
	return &transferCoin{prompterApp: prompterApp, massaClient: massaClient}
}

type transferCoin struct {
	prompterApp prompt.WalletPrompterInterface
	massaClient network.NodeFetcherInterface
}

func (t *transferCoin) Handle(params operations.TransferCoinParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	wlt, resp := loadWallet(params.Nickname)
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

	promptData := &prompt.PromptRequestData{
		Msg:  fmt.Sprintf("Transfer %s nonaMassa from %s to %s, with fee %s nonaMassa", string(params.Body.Amount), wlt.Nickname, *params.Body.RecipientAddress, string(params.Body.Fee)),
		Data: nil,
	}

	_, err = prompt.PromptPassword(t.prompterApp, wlt, walletapp.Transfer, promptData)
	if err != nil {
		return operations.NewTransferCoinUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	// create the transaction and send it to the network
	operation, err := doTransfer(wlt, amount, fee, *params.Body.RecipientAddress, t.massaClient)
	if err != nil {
		errStr := fmt.Sprintf("error transferring coin: %v", err.Error())
		t.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
			walletapp.EventData{Success: false, Data: errStr})
		return operations.NewTransferCoinInternalServerError().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: errStr,
			})
	}

	t.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
		walletapp.EventData{Success: true, Data: "Transfer Success"})
	return operations.NewTransferCoinOK().WithPayload(
		&models.OperationResponse{
			OperationID: operation.OperationID,
		})
}

func doTransfer(wlt *wallet.Wallet, amount, fee uint64, recipientAddress string, massaClient network.NodeFetcherInterface) (*sendOperation.OperationResponse, error) {
	operation, err := transaction.New(recipientAddress, amount)
	if err != nil {
		return nil, fmt.Errorf("Error during transaction creation: %w", err)
	}

	msg, err := massaClient.MakeOperation(fee, operation)
	if err != nil {
		return nil, fmt.Errorf("Error while making operation: %w", err)
	}

	// sign the msg in base64
	// TODO: we do not implement the handling of the correlation id for now
	byteMsgB64 := strfmt.Base64(msg)
	signature, err := wlt.Sign(&byteMsgB64)
	if err != nil {
		return nil, fmt.Errorf("Error sign: %w", err)
	}

	// send the msg to the network

	resp, err := massaClient.MakeRPCCall(msg, signature, wlt.GetPupKey())
	if err != nil {
		return nil, fmt.Errorf("Error during RPC call: %w", err)
	}

	return &sendOperation.OperationResponse{CorrelationID: "", OperationID: resp[0]}, nil
}

// loadWalletToTransfer loads a wallet from the file system or returns an error.
func loadWalletToTransfer(nickname string) (*wallet.Wallet, middleware.Responder) {
	w, err := wallet.Load(nickname)
	if err == nil {
		return w, nil
	}

	errorObj := models.Error{
		Code:    errorGetWallets,
		Message: err.Error(),
	}

	if err.Error() == wallet.ErrorAccountNotFound(nickname).Error() {
		return nil, operations.NewTransferCoinNotFound().WithPayload(&errorObj)
	} else {
		return nil, operations.NewTransferCoinBadRequest().WithPayload(&errorObj)
	}
}
