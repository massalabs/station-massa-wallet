package wallet

import (
	"bytes"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/awnumar/memguard"
	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station-massa-wallet/api/server/models"
	"github.com/massalabs/station-massa-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
	"github.com/pkg/errors"
	"lukechampine.com/blake3"
)

const (
	passwordExpirationTime = time.Second * 60 * 10
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

func (w *walletSign) Handle(params operations.SignParams) middleware.Responder {
	acc, resp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if resp != nil {
		return resp
	}

	promptRequest, err := w.getPromptRequest(params.Body.Operation.String(), acc, params.Body.Description)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %s", errorSignDecodeMessage), errorSignDecodeMessage, http.StatusBadRequest)
	}

	var correlationID models.CorrelationID
	var privateKey *memguard.LockedBuffer

	if params.Body.CorrelationID != nil {
		correlationID = params.Body.CorrelationID

		pk, err := w.privateKeyFromCache(acc, correlationID)
		if err != nil {
			if errors.Is(err, utils.ErrCorrelationIDNotFound) {
				return newErrorResponse(err.Error(), errorSignCorrelationIdNotFound, http.StatusNotFound)
			}

			return newErrorResponse(err.Error(), errorSign, http.StatusInternalServerError)
		}

		privateKey = pk
	} else {
		password, err := w.PromptPassword(acc, promptRequest)
		if err != nil {
			msg := fmt.Sprintf("Unable to unprotect wallet: %s", err.Error())
			if errors.Is(err, utils.ErrWrongPassword) || errors.Is(err, utils.ErrActionCanceled) {
				return newErrorResponse(msg, errorGetWallets, http.StatusUnauthorized)
			}

			return newErrorResponse(msg, errorGetWallets, http.StatusInternalServerError)
		}

		pk, err := acc.PrivateKeyBytesInClear(password)
		if err != nil {
			return newErrorResponse(err.Error(), errorWrongPassword, http.StatusInternalServerError)
		}

		privateKey = pk

		if params.Body.Batch {
			cID, err := w.CacheAccount(acc, privateKey)
			if err != nil {
				return newErrorResponse(err.Error(), errorSignCorrelationIdNotFound, http.StatusInternalServerError)
			}
			correlationID = *cID
		}
	}

	operation, err := w.prepareOperation(acc, params, promptRequest)
	if err != nil {
		return newErrorResponse(err.Error(), errorSignDecodeOperation, http.StatusBadRequest)
	}

	signature := acc.SignWithPrivateKey(privateKey, operation)

	return w.Success(acc, signature, correlationID)
}

func (w *walletSign) PromptPassword(acc *account.Account, promptRequest *prompt.PromptRequest) (*memguard.LockedBuffer, error) {
	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, *promptRequest, acc)
	if err != nil {
		return nil, fmt.Errorf("prompting password: %w", err)
	}

	password, _ := promptOutput.(*memguard.LockedBuffer)

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	return password, nil
}

func (w *walletSign) CacheAccount(acc *account.Account, privateKey *memguard.LockedBuffer) (*models.CorrelationID, error) {
	correlationID, err := generateCorrelationId()
	if err != nil {
		return nil, fmt.Errorf("Error cannot generate correlation id: %w", err)
	}

	key := CacheKey(correlationID)

	cacheValue, err := Xor(privateKey.Bytes(), correlationID)
	if err != nil {
		return nil, fmt.Errorf("Error cannot XOR correlation id: %w", err)
	}

	err = w.gc.SetWithExpire(key, cacheValue, expirationDuration())
	if err != nil {
		return nil, fmt.Errorf("Error set correlation id in cache: %v", err.Error())
	}

	return &correlationID, nil
}

func (w *walletSign) Success(acc *account.Account, signature []byte, correlationId models.CorrelationID) middleware.Responder {
	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	return operations.NewSignOK().WithPayload(
		&models.SignResponse{
			PublicKey:     string(publicKeyBytes),
			Signature:     signature,
			CorrelationID: correlationId,
		})
}

func (w *walletSign) prepareOperation(acc *account.Account, params operations.SignParams, promptRequest *prompt.PromptRequest) ([]byte, error) {
	operation, err := base64.StdEncoding.DecodeString(params.Body.Operation.String())
	if err != nil {
		return nil, fmt.Errorf("Unable to decode operation: %w", err)
	}

	if promptRequest.Data.(PromptRequestSignData).OperationType != Message {
		publicKeyBytes, err := acc.PublicKey.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("Unable to marshal public key: %w", err)
		}

		operation = append(publicKeyBytes, operation...)
	}

	return operation, nil
}

func (w *walletSign) getPromptRequest(msgToSign string, acc *account.Account, description string) (*prompt.PromptRequest, error) {
	var promptRequest prompt.PromptRequest
	var opType uint64
	var err error

	addressBytes, err := acc.Address.MarshalText()
	if err != nil {
		return nil, err
	}
	address := string(addressBytes)

	decodedMsg, fees, expiry, err := sendoperation.DecodeMessage64(msgToSign)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode transaction message")
	}

	if opType, err = sendoperation.DecodeOperationID(decodedMsg); err != nil {
		wrappedErr := errors.Wrap(err, "failed to decode operation ID")

		return nil, wrappedErr
	} else {
		switch opType {
		case TransactionOpType:
			msg, err := transaction.DecodeMessage(decodedMsg)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode transaction message")
			}
			promptRequest = w.prepareTransactionPromptRequest(msg, acc, address, description, fees, expiry)

		case BuyRollOpType, SellRollOpType:
			roll, err := sendoperation.RollDecodeMessage(decodedMsg)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode roll message")
			}
			promptRequest = w.prepareRollPromptRequest(roll, acc, address, description, fees, expiry)

		case ExecuteSCOpType:
			executeSC, err := executesc.DecodeMessage(decodedMsg)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode executeSC message")
			}
			promptRequest = w.prepareExecuteSCPromptRequest(executeSC, acc, address, description, fees, expiry)

		case CallSCOpType:
			callSC, err := callsc.DecodeMessage(decodedMsg)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode callSC message")
			}
			promptRequest = w.prepareCallSCPromptRequest(callSC, acc, address, description, fees, expiry)

		default:
			decodedMsg, err := base64.StdEncoding.DecodeString(msgToSign)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode plainText message from b64")
			}
			promptRequest = w.prepareplainTextPromptRequest(string(decodedMsg), acc, address, description)

		}
	}

	return &promptRequest, nil
}

func (w *walletSign) prepareCallSCPromptRequest(msg *callsc.MessageContent,
	acc *account.Account,
	address string,
	description string,
	fees uint64,
	expiry uint64,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", acc.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			Fees:          strconv.FormatUint(fees, 10),
			OperationType: "Call SC",
			OperationID:   msg.OperationID,
			MaxGas:        strconv.FormatUint(msg.MaxGas, 10),
			Coins:         strconv.FormatUint(msg.Coins, 10),
			Address:       msg.Address,
			Function:      msg.Function,
			Expiry:        expiry,
			WalletAddress: address,
			Nickname:      acc.Nickname,
		},
	}
}

func (w *walletSign) prepareExecuteSCPromptRequest(
	msg *executesc.MessageContent,
	acc *account.Account,
	address string,
	description string,
	fees uint64,
	expiry uint64,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", acc.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			Fees:          strconv.FormatUint(fees, 10),
			OperationType: "Execute SC",
			MaxCoins:      strconv.FormatUint(msg.MaxCoins, 10),
			MaxGas:        strconv.FormatUint(msg.MaxGas, 10),
			Expiry:        expiry,
			WalletAddress: address,
			Nickname:      acc.Nickname,
		},
	}
}

func (w *walletSign) prepareRollPromptRequest(
	msg *sendoperation.RollMessageContent,
	acc *account.Account,
	address string,
	description string,
	fees uint64,
	expiry uint64,
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
		Msg:    fmt.Sprintf("Unprotect wallet %s", acc.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			Fees:          strconv.FormatUint(fees, 10),
			OperationType: operationType,
			RollCount:     msg.RollCount,
			Expiry:        expiry,
			WalletAddress: address,
			Nickname:      acc.Nickname,
		},
	}
}

func (w *walletSign) prepareTransactionPromptRequest(
	msg *transaction.MessageContent,
	acc *account.Account,
	address string,
	description string,
	fees uint64,
	expiry uint64,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", acc.Nickname),
		Data: PromptRequestSignData{
			Description:      description,
			Fees:             strconv.FormatUint(fees, 10),
			OperationType:    "Transaction",
			RecipientAddress: msg.RecipientAddress,
			Amount:           strconv.FormatUint(msg.Amount, 10),
			Expiry:           expiry,
			WalletAddress:    address,
			Nickname:         acc.Nickname,
		},
	}
}

func (s *walletSign) prepareplainTextPromptRequest(
	plainText string,
	acc *account.Account,
	address string,
	description string,
) prompt.PromptRequest {
	return prompt.PromptRequest{
		Action: walletapp.Sign,
		Msg:    fmt.Sprintf("Unprotect wallet %s", acc.Nickname),
		Data: PromptRequestSignData{
			Description:   description,
			OperationType: Message,
			PlainText:     plainText,
			Nickname:      acc.Nickname,
			WalletAddress: address,
		},
	}
}

// privateKeyFromCache return the private key from the cache or an error.
func (w *walletSign) privateKeyFromCache(
	acc *account.Account,
	correlationID models.CorrelationID,
) (*memguard.LockedBuffer, error) {
	key := CacheKey(correlationID)

	value, err := w.gc.Get(key)
	if err != nil {
		if err.Error() == gcache.KeyNotFoundError.Error() {
			return nil, fmt.Errorf("%w: %w", utils.ErrCorrelationIDNotFound, err)
		}

		return nil, fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	// convert interface{} into byte[]
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	cipheredPrivateKey := buf.Bytes()

	privateKey, err := Xor(cipheredPrivateKey, correlationID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	return memguard.NewBufferFromBytes(privateKey), nil
}

func Xor(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of two arrays must be same, %d and %d", len(a), len(b))
	}
	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}

	return result, nil
}

func CacheKey(correlationID models.CorrelationID) [32]byte {
	return blake3.Sum256(correlationID)
}

func expirationDuration() time.Duration {
	fromEnv := os.Getenv("BATCH_EXPIRATION_TIME")

	if fromEnv == "" {
		return passwordExpirationTime
	}

	duration, err := time.ParseDuration(fromEnv)
	if err != nil {
		return passwordExpirationTime
	}

	return duration
}

func generateCorrelationId() (models.CorrelationID, error) {
	rand := cryptorand.Reader

	correlationId := make([]byte, ed25519.PrivateKeySize)
	if _, err := io.ReadFull(rand, correlationId); err != nil {
		return nil, err
	}

	return correlationId, nil
}
