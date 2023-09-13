package wallet

import (
	"bytes"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
	"github.com/pkg/errors"
	"lukechampine.com/blake3"
)

const (
	passwordExpirationTime = time.Second * 60 * 30
	BuyRoll                = "Buy Roll"
	SellRoll               = "Sell Roll"
	Message                = "Plain Text"
	TransactionOpType      = uint64(0)
	BuyRollOpType          = uint64(1)
	SellRollOpType         = uint64(2)
	ExecuteSCOpType        = uint64(3)
	CallSCOpType           = uint64(4)
)

type PromptRequestSignData struct {
	Description      string
	Fees             string
	OperationType    string
	OperationID      uint64
	GasLimit         string
	Coins            string
	Address          string
	Function         string
	MaxCoins         string
	MaxGas           string
	Expiry           uint64
	WalletAddress    string
	Nickname         string
	RollCount        uint64
	RecipientAddress string
	Amount           string
	PlainText        string
}

// NewSign instantiates a sign Handler
// The "classical" way is not possible because we need to pass to the handler a password.PasswordAsker.
func NewSign(prompterApp prompt.WalletPrompterInterface, gc gcache.Cache) operations.SignHandler {
	return &walletSign{gc: gc, prompterApp: prompterApp}
}

type walletSign struct {
	prompterApp prompt.WalletPrompterInterface
	gc          gcache.Cache
}

func (s *walletSign) Handle(params operations.SignParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	var correlationId models.CorrelationID
	promptRequest, err := s.getPromptRequest(params.Body.Operation.String(), wlt, params.Body.Description)
	if err != nil {
		return s.handleBadRequest(errorSignDecodeMessage)
	}

	if params.Body.CorrelationID != nil {
		resp = handleWithCorrelationId(wlt, params, s.gc)
		correlationId = params.Body.CorrelationID
	} else {
		_, err = prompt.WakeUpPrompt(s.prompterApp, promptRequest, wlt)
		if err != nil {
			return operations.NewSignUnauthorized().WithPayload(
				&models.Error{
					Code:    errorCanceledAction,
					Message: "Unable to unprotect wallet",
				})
		}

		s.prompterApp.EmitEvent(walletapp.PromptResultEvent,
			walletapp.EventData{Success: true, CodeMessage: utils.MsgAccountUnprotected})

		if params.Body.Batch {
			correlationId, resp = handleBatch(wlt, params, s, s.gc)
		}
	}

	if resp != nil {
		return resp
	}

	signingOperation := true

	if promptRequest.Data.(PromptRequestSignData).OperationType == Message {
		signingOperation = false
	}

	op, err := base64.StdEncoding.DecodeString(params.Body.Operation.String())
	if err != nil {
		return s.handleBadRequest(errorSignDecodeOperation)
	}
	signature := wlt.Sign(signingOperation, op)

	return operations.NewSignOK().WithPayload(
		&models.SignResponse{
			PublicKey:     wlt.GetPupKey(),
			Signature:     signature,
			CorrelationID: correlationId,
		})
}

func (s *walletSign) handleBadRequest(errorCode string) middleware.Responder {
	s.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: false, CodeMessage: errorCode})
	return operations.NewSignBadRequest().WithPayload(
		&models.Error{
			Code:    errorCode,
			Message: fmt.Sprintf("Error: %s", errorCode),
		})
}

func (s *walletSign) getPromptRequest(msgToSign string, wlt *wallet.Wallet, description string) (prompt.PromptRequest, error) {
	var promptRequest prompt.PromptRequest
	var opType uint64
	var err error

	decodedMsg, fees, expiry, err := sendoperation.DecodeMessage64(msgToSign)
	if err != nil {
		return s.prepareUnknownPromptRequest(wlt, description), errors.Wrap(err, "failed to decode transaction message")
	}

	if opType, err = sendoperation.DecodeOperationID(decodedMsg); err != nil {
		wrappedErr := errors.Wrap(err, "failed to decode operation ID")

		return s.prepareUnknownPromptRequest(wlt, description), wrappedErr
	} else {
		switch opType {
		case TransactionOpType:
			msg, err := transaction.DecodeMessage(decodedMsg)
			if err != nil {
				wrappedErr := errors.Wrap(err, "failed to decode transaction message")
				return s.prepareUnknownPromptRequest(wlt, description), wrappedErr
			}
			promptRequest = s.prepareTransactionPromptRequest(msg, wlt, description, fees)

		case BuyRollOpType, SellRollOpType:
			roll, err := sendoperation.RollDecodeMessage(decodedMsg)
			if err != nil {
				wrappedErr := errors.Wrap(err, "failed to decode roll message")
				return s.prepareUnknownPromptRequest(wlt, description), wrappedErr
			}
			promptRequest = s.prepareRollPromptRequest(roll, wlt, description, fees)

		case ExecuteSCOpType:
			executeSC, err := executesc.DecodeMessage(decodedMsg)
			if err != nil {
				wrappedErr := errors.Wrap(err, "failed to decode executeSC message")
				return s.prepareUnknownPromptRequest(wlt, description), wrappedErr
			}
			promptRequest = s.prepareExecuteSCPromptRequest(executeSC, wlt, description, fees, expiry)

		case CallSCOpType:
			callSC, err := callsc.DecodeMessage(decodedMsg)
			if err != nil {
				wrappedErr := errors.Wrap(err, "failed to decode callSC message")
				return s.prepareUnknownPromptRequest(wlt, description), wrappedErr
			}
			promptRequest = s.prepareCallSCPromptRequest(callSC, wlt, description, fees)

		default:
			decodedMsg, err := base64.StdEncoding.DecodeString(msgToSign)
			if err != nil {
				wrappedErr := errors.Wrap(err, "failed to decode plainText message from b64")
				return s.prepareUnknownPromptRequest(wlt, description), wrappedErr
			}
			promptRequest = s.prepareplainTextPromptRequest(string(decodedMsg), wlt, description)
		}
	}

	return promptRequest, nil
}

func (s *walletSign) prepareCallSCPromptRequest(msg *callsc.MessageContent,
	wlt *wallet.Wallet,
	description string,
	fees uint64,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			Fees:          strconv.FormatUint(fees, 10),
			OperationType: "Call SC",
			OperationID:   msg.OperationID,
			GasLimit:      strconv.FormatUint(msg.GasLimit, 10),
			Coins:         strconv.FormatUint(msg.Coins, 10),
			Address:       msg.Address,
			Function:      msg.Function,
			WalletAddress: wlt.Address,
			Nickname:      wlt.Nickname,
		},
	}
}

func (s *walletSign) prepareExecuteSCPromptRequest(
	msg *executesc.MessageContent,
	wlt *wallet.Wallet,
	description string,
	fees uint64,
	expiry uint64,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			Fees:          strconv.FormatUint(fees, 10),
			OperationType: "Execute SC",
			MaxCoins:      strconv.FormatUint(msg.MaxCoins, 10),
			MaxGas:        strconv.FormatUint(msg.MaxGas, 10),
			Expiry:        expiry,
			WalletAddress: wlt.Address,
			Nickname:      wlt.Nickname,
		},
	}
}

func (s *walletSign) prepareRollPromptRequest(
	msg *sendoperation.RollMessageContent,
	wlt *wallet.Wallet,
	description string,
	fees uint64,
) prompt.PromptRequest {
	operationType := ""
	switch msg.OperationID {
	case 1:
		operationType = BuyRoll
	case 2:
		operationType = SellRoll
	}

	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			Fees:          strconv.FormatUint(fees, 10),
			OperationType: operationType,
			RollCount:     msg.RollCount,
			WalletAddress: wlt.Address,
			Nickname:      wlt.Nickname,
		},
	}
}

func (s *walletSign) prepareTransactionPromptRequest(
	msg *transaction.MessageContent,
	wlt *wallet.Wallet,
	description string,
	fees uint64,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
		Data: PromptRequestSignData{
			Description:      description,
			Fees:             strconv.FormatUint(fees, 10),
			OperationType:    "Transaction",
			RecipientAddress: msg.RecipientAddress,
			Amount:           strconv.FormatUint(msg.Amount, 10),
			WalletAddress:    wlt.Address,
			Nickname:         wlt.Nickname,
		},
	}
}

func (s *walletSign) prepareplainTextPromptRequest(
	plainText string,
	wlt *wallet.Wallet,
	description string,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			OperationType: Message,
			PlainText:     plainText,
			WalletAddress: wlt.Address,
			Nickname:      wlt.Nickname,
		},
	}
}

func (s *walletSign) prepareUnknownPromptRequest(wlt *wallet.Wallet, description string) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			OperationType: Message,
			WalletAddress: wlt.Address,
			Nickname:      wlt.Nickname,
		},
	}
}

func handleWithCorrelationId(
	wlt *wallet.Wallet,
	params operations.SignParams,
	gc gcache.Cache,
) middleware.Responder {
	key := CacheKey(params.Body.CorrelationID)

	value, err := gc.Get(key)
	if err != nil {
		if err == gcache.KeyNotFoundError {
			return operations.NewSignNotFound().WithPayload(
				&models.Error{
					Code:    errorSignCorrelationIdNotFound,
					Message: fmt.Sprintf("Error given correlation id not in cache: %v", err.Error()),
				})
		}
		return operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: fmt.Sprintf("Error cannot get data from cache: %v", err.Error()),
			})
	}

	// convert interface{} into byte[]
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: fmt.Sprintf("Error cannot convert cache value: %v", err.Error()),
			})
	}

	bytes := buf.Bytes()

	err = wlt.UnprotectFromCorrelationId(bytes, params.Body.CorrelationID)
	if err != nil {
		return operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: fmt.Sprintf("Error cannot unprotect from cache: %v", err.Error()),
			})
	}

	return nil
}

func CacheKey(correlationId models.CorrelationID) [32]byte {
	return blake3.Sum256(correlationId)
}

func handleBatch(
	wlt *wallet.Wallet,
	params operations.SignParams,
	s *walletSign,
	gc gcache.Cache,
) (models.CorrelationID, middleware.Responder) {
	correlationId, err := generateCorrelationId()
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignGenerateCorrelationId,
				Message: fmt.Sprintf("Error cannot generate correlation id: %v", err.Error()),
			})
	}

	key := CacheKey(correlationId)

	cacheValue, err := wallet.Xor(wlt.KeyPair.PrivateKey, correlationId)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignGenerateCorrelationId,
				Message: fmt.Sprintf("Error cannot XOR correlation id: %v", err.Error()),
			})
	}

	err = gc.SetWithExpire(key, cacheValue, passwordExpirationTime)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignGenerateCorrelationId,
				Message: fmt.Sprintf("Error set correlation id in cache: %v", err.Error()),
			})
	}

	return correlationId, nil
}

func generateCorrelationId() (models.CorrelationID, error) {
	rand := cryptorand.Reader

	// Correlation id must have the same size as the versioned private key.
	correlationId := make([]byte, ed25519.PrivateKeySize+1)
	if _, err := io.ReadFull(rand, correlationId); err != nil {
		return nil, err
	}

	return correlationId, nil
}
