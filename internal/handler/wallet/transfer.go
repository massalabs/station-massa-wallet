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
	"github.com/massalabs/thyra/pkg/node"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/transaction"
)

func NewTransferCoin(prompterApp prompt.WalletPrompterInterface) operations.TransferCoinHandler {
	return &wTransferCoin{prompterApp: prompterApp}
}

type wTransferCoin struct {
	prompterApp prompt.WalletPrompterInterface
}

func (h *wTransferCoin) Handle(params operations.TransferCoinParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	wlt, resp := loadWalletForTransfer(params.Nickname)
	if resp != nil {
		return resp
	}

	promptData := &prompt.PromptRequestData{
		Msg:  fmt.Sprintf("Transfer %s nonaMassa from %s to %s, with fee %s nonaMassa", string(params.Body.Amount), wlt.Nickname, *params.Body.RecipientAddress, string(params.Body.Fee)),
		Data: nil,
	}

	_, err := prompt.PromptPassword(h.prompterApp, wlt, walletapp.Transfer, promptData)
	if err != nil {
		return operations.NewTransferCoinUnauthorized().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: "Unable to unprotect wallet",
			})
	}

	// create the transaction and send it to the network
	operation, err := doTransfer(wlt, params.Body)
	if err != nil {
		errStr := "error transferring coin:" + err.Error()
		h.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
			walletapp.EventData{Success: false, Data: errStr})
		return operations.NewTransferCoinInternalServerError().WithPayload(
			&models.Error{
				Code:    errorTransferCoin,
				Message: errStr,
			})
	}

	h.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
		walletapp.EventData{Success: true, Data: "Transfer Success"})
	return operations.NewTransferCoinOK().WithPayload(
		&models.OperationResponse{
			OperationID: &operation.OperationID,
		})
}

func doTransfer(wlt *wallet.Wallet, body *models.TransferRequest) (*sendOperation.OperationResponse, error) {
	recipientAddress := *body.RecipientAddress
	netWorkInfo, err := network.GetNetworkInfo()
	if err != nil {
		return nil, fmt.Errorf("Error during network info retrieval: %w", err)
	}

	url := netWorkInfo.URL
	client := node.NewClient(url)

	// convert amount to uint64
	amount, err := strconv.ParseUint(string(body.Amount), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error during conversion")
	}

	// convert fee to uint64
	fee, err := strconv.ParseUint(string(body.Fee), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error during conversion")
	}

	operation, err := transaction.New(recipientAddress, amount)
	if err != nil {
		return nil, fmt.Errorf("Error during transaction creation: %w", err)
	}
	msg, _, err := sendOperation.MakeOperation(client, sendOperation.DefaultSlotsDuration, fee, operation)
	if err != nil {
		return nil, fmt.Errorf("Error during operation creation: %w", err)
	}

	// sign the msg in base64
	// TODO: we do not implement the handling of the correlation id for now
	byteMsgB64 := strfmt.Base64(msg)
	signature, err := wlt.Sign(&byteMsgB64)
	if err != nil {
		return nil, fmt.Errorf("Error sign: %w", err)
	}

	// send the msg to the network
	resp, err := sendOperation.MakeRPCCall(msg, signature, wlt.GetPupKey(), client)
	if err != nil {
		return nil, fmt.Errorf("Error during RPC call: %w", err)
	}

	return &sendOperation.OperationResponse{CorrelationID: "", OperationID: resp[0]}, nil
}

// loadWalletForSign loads a wallet from the file system or returns an error.
func loadWalletForTransfer(nickname string) (*wallet.Wallet, middleware.Responder) {
	w, err := wallet.Load(nickname)
	if err != nil {
		if err.Error() == wallet.ErrorAccountNotFound(nickname).Error() {
			return nil, operations.NewTransferCoinNotFound().WithPayload(
				&models.Error{
					Code:    errorGetWallets,
					Message: err.Error(),
				})
		} else {
			return nil, operations.NewTransferCoinBadRequest().WithPayload(
				&models.Error{
					Code:    errorGetWallets,
					Message: err.Error(),
				})
		}
	}

	return w, nil
}
