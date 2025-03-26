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
	"github.com/massalabs/station-massa-wallet/pkg/cache"
	"github.com/massalabs/station-massa-wallet/pkg/config"
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
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
	acc, errResp := loadAccount(t.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	// convert amount to uint64
	amount, err := strconv.ParseUint(string(params.Body.Amount), 10, 64)
	if err != nil {
		return newErrorResponse("Error during amount conversion", errorTransferCoin, http.StatusBadRequest)
	}

	// convert fee to uint64
	fee, err := strconv.ParseUint(string(params.Body.Fee), 10, 64)
	if err != nil {
		return newErrorResponse("Error during fee conversion", errorTransferCoin, http.StatusBadRequest)
	}

	addressBytes, err := acc.Address.MarshalText()
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}
	address := string(addressBytes)

	var recipientNickname string

	recipientAcc, err := t.prompterApp.App().Wallet.GetAccountFromAddress(*params.Body.RecipientAddress)
	if err != nil {
		recipientNickname = ""
	} else {
		recipientNickname = recipientAcc.Nickname
	}

	chainID, minimalFees, err := network.GetNodeInfo()
	if err != nil {
		return newErrorResponse("failed to get chain id", errorTransferCoin, http.StatusInternalServerError)
	}

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Sign,
		Data: prompt.PromptRequestSignData{
			Fees:              strconv.FormatUint(fee, 10),
			MinFees:           minimalFees,
			WalletAddress:     address,
			Nickname:          acc.Nickname,
			OperationType:     int(transaction.OpType),
			AllowFeeEdition:   true,
			RecipientAddress:  *params.Body.RecipientAddress,
			RecipientNickname: recipientNickname,
			Amount:            string(params.Body.Amount),
			ChainID:           int64(chainID),
		},
	}

	promptOutput, err := prompt.WakeUpPrompt(t.prompterApp, promptRequest, acc)
	if err != nil {
		return newErrorResponse("Unable to unprotect wallet", errorCanceledAction, http.StatusUnauthorized)
	}

	output, ok := promptOutput.(*walletapp.SignPromptOutput)
	if !ok {
		return newErrorResponse(fmt.Sprintf("prompting password for message: %v", utils.ErrInvalidInputType.Error()), utils.ErrInvalidInputType.Error(), http.StatusInternalServerError)
	}
	password := output.Password

	// create the transaction and send it to the network
	operation, err := doTransfer(acc, password, amount, output.Fees, *params.Body.RecipientAddress, t.massaClient, chainID)
	if err != nil {
		msg := fmt.Sprintf("error transferring coin: %v", err.Error())

		t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: false, CodeMessage: errorTransferCoin})

		return newErrorResponse(msg, errorTransferCoin, http.StatusInternalServerError)
	}

	t.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	cfg := config.Get()

	if cfg.HasEnabledRule(acc.Nickname) {
		err = cache.CachePrivateKeyFromPassword(acc, output.Password)
		if err != nil {
			return newErrorResponse(err.Error(), errorCachePrivateKey, http.StatusInternalServerError)
		}
	}

	return operations.NewTransferCoinOK().WithPayload(
		&models.OperationResponse{
			OperationID: operation.OperationID,
		})
}

func doTransfer(
	acc *account.Account,
	password *memguard.LockedBuffer,
	amount, fee uint64,
	recipientAddress string,
	massaClient network.NodeFetcherInterface,
	chainID uint64,
) (*sendOperation.OperationResponse, error) {
	operation, err := transaction.New(recipientAddress, amount)
	if err != nil {
		return nil, fmt.Errorf("error during transaction creation: %w", err)
	}

	return network.SendOperation(acc, password, massaClient, operation, fee, chainID)
}
