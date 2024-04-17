package wallet

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
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
	"github.com/massalabs/station-massa-wallet/pkg/network"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/utils"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/buyrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/callsc"
	"github.com/massalabs/station/pkg/node/sendoperation/executesc"
	"github.com/massalabs/station/pkg/node/sendoperation/sellrolls"
	"github.com/massalabs/station/pkg/node/sendoperation/transaction"
	"github.com/pkg/errors"
	"lukechampine.com/blake3"
)

const (
	defaultExpirationTime = time.Second * 60 * 10
	RollPrice             = 100
)

type PromptRequestSignData struct {
	Description       string
	Fees              string
	MinFees           string
	OperationType     int
	Coins             string
	Address           string
	Function          string
	MaxCoins          string
	WalletAddress     string
	Nickname          string
	RollCount         uint64
	RecipientAddress  string
	RecipientNickname string
	Amount            string
	PlainText         string
	AllowFeeEdition   bool
	ChainID           int64
}

func NewSign(prompterApp prompt.WalletPrompterInterface, gc gcache.Cache) operations.SignHandler {
	return &walletSign{gc: gc, prompterApp: prompterApp}
}

type walletSign struct {
	prompterApp prompt.WalletPrompterInterface
	gc          gcache.Cache
}

func (w *walletSign) Handle(params operations.SignParams) middleware.Responder {
	acc, errResp := loadAccount(w.prompterApp.App().Wallet, params.Nickname)
	if errResp != nil {
		return errResp
	}

	promptRequest, fees, err := w.getPromptRequest(params, acc, params.Body.Description, *params.Body.ChainID)
	if err != nil {
		return newErrorResponse(fmt.Sprintf("Error: %v", err.Error()), errorSignDecodeMessage, http.StatusBadRequest)
	}

	var correlationID *memguard.LockedBuffer
	var privateKey *memguard.LockedBuffer

	const disableBatchSigning = true

	if params.Body.CorrelationID != nil && !disableBatchSigning {
		correlationID = memguard.NewBufferFromBytes(params.Body.CorrelationID)

		pk, err := w.privateKeyFromCache(acc, correlationID)
		if err != nil {
			if errors.Is(err, utils.ErrCorrelationIDNotFound) {
				return newErrorResponse(err.Error(), errorSignCorrelationIDNotFound, http.StatusNotFound)
			}

			return newErrorResponse(err.Error(), errorSign, http.StatusInternalServerError)
		}

		privateKey = pk
	} else {
		output, err := w.PromptPassword(acc, promptRequest)
		if err != nil {
			msg := fmt.Sprintf("Unable to unprotect wallet: %s", err.Error())
			if errors.Is(err, utils.ErrWrongPassword) || errors.Is(err, utils.ErrActionCanceled) {
				return newErrorResponse(msg, errorGetWallets, http.StatusUnauthorized)
			}

			return newErrorResponse(msg, errorGetWallets, http.StatusInternalServerError)
		}

		fees = output.Fees

		pk, err := acc.PrivateKeyBytesInClear(output.Password)
		if err != nil {
			return newErrorResponse(err.Error(), errorWrongPassword, http.StatusInternalServerError)
		}

		privateKey = pk

		if params.Body.Batch {
			cID, err := w.CacheAccount(acc, privateKey)
			if err != nil {
				return newErrorResponse(err.Error(), errorSignCorrelationIDNotFound, http.StatusInternalServerError)
			}
			correlationID = cID
		} else {
			correlationID = memguard.NewBufferFromBytes([]byte{})
		}
	}

	operation, msgToSign, err := prepareOperation(acc, fees, params.Body.Operation.String(), promptRequest.Data.(PromptRequestSignData).OperationType, *params.Body.ChainID)
	if err != nil {
		return newErrorResponse(err.Error(), errorSignDecodeOperation, http.StatusBadRequest)
	}

	signature := acc.SignWithPrivateKey(privateKey, msgToSign)

	if !acc.VerifySignature(msgToSign, signature) {
		return newErrorResponse("Error: signature verification failed", "errorSignVerifySignature", http.StatusInternalServerError)
	}

	return w.Success(acc, signature, correlationID, operation)
}

func (w *walletSign) PromptPassword(acc *account.Account, promptRequest *prompt.PromptRequest) (*walletapp.SignPromptOutput, error) {
	promptOutput, err := prompt.WakeUpPrompt(w.prompterApp, *promptRequest, acc)
	if err != nil {
		return nil, fmt.Errorf("prompting password: %w", err)
	}

	output, ok := promptOutput.(*walletapp.SignPromptOutput)
	if !ok {
		return nil, fmt.Errorf("prompting password for sign: %s", utils.ErrInvalidInputType.Error())
	}

	w.prompterApp.EmitEvent(walletapp.PromptResultEvent,
		walletapp.EventData{Success: true})

	return output, nil
}

func (w *walletSign) CacheAccount(acc *account.Account, privateKey *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	correlationID, err := generateCorrelationID()
	if err != nil {
		return nil, fmt.Errorf("Error cannot generate correlation id: %w", err)
	}

	key := CacheKey(correlationID.Bytes())

	cipheredPrivateKey, err := Xor(privateKey, correlationID)
	if err != nil {
		return nil, fmt.Errorf("Error cannot XOR correlation id: %w", err)
	}

	cacheValue := make([]byte, ed25519.PrivateKeySize)
	copy(cacheValue, cipheredPrivateKey.Bytes())
	cipheredPrivateKey.Destroy()

	err = w.gc.SetWithExpire(key, cacheValue, expirationDuration())
	if err != nil {
		return nil, fmt.Errorf("Error set correlation id in cache: %w", err)
	}

	return correlationID, nil
}

func (w *walletSign) Success(acc *account.Account, signature []byte, correlationId *memguard.LockedBuffer, operation []byte) middleware.Responder {
	publicKeyBytes, err := acc.PublicKey.MarshalText()
	if err != nil {
		return newErrorResponse(err.Error(), errorGetAccount, http.StatusInternalServerError)
	}

	return operations.NewSignOK().WithPayload(
		&models.SignResponse{
			PublicKey:     string(publicKeyBytes),
			Signature:     signature,
			CorrelationID: correlationId.Bytes(),
			Operation:     operation,
		})
}

// prepareOperation prepares the operation to be signed.
// Returns the modified operation (fees change) and the operation to be signed (with public key).
// Returns an error if the operation cannot be decoded.
func prepareOperation(acc *account.Account, fees uint64, operationB64 string, operationType int, chainID int64) ([]byte, []byte, error) {
	decodedMsg, _, expiry, err := sendoperation.DecodeMessage64(operationB64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode operation for preparing before signing: %w", err)
	}

	operation := make([]byte, 0)
	operation = binary.AppendUvarint(operation, fees)
	operation = binary.AppendUvarint(operation, expiry)
	operation = append(operation, decodedMsg...)

	// operation to be signed
	msgToSign := make([]byte, len(operation))
	copy(msgToSign, operation)

	publicKey, err := acc.PublicKey.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to marshal public key: %w", err)
	}

	msgToSign = utils.PrepareSignData(uint64(chainID), append(publicKey, msgToSign...))

	return operation, msgToSign, nil
}

func (w *walletSign) getPromptRequest(params operations.SignParams, acc *account.Account, description string, chainID int64) (*prompt.PromptRequest, uint64, error) {
	msgToSign := params.Body.Operation.String()

	decodedMsg, fees, _, err := sendoperation.DecodeMessage64(msgToSign)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode transaction message: %w", err)
	}

	opType, err := sendoperation.DecodeOperationType(decodedMsg)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode operation ID: %w", err)
	}

	var data PromptRequestSignData

	switch opType {
	case transaction.OpType:
		data, err = w.getTransactionPromptData(decodedMsg, acc)

	case buyrolls.OpType, sellrolls.OpType:
		data, err = getRollPromptData(decodedMsg, acc)

	case executesc.OpType:
		data, err = getExecuteSCPromptData(decodedMsg, acc)

	case callsc.OpType:
		data, err = getCallSCPromptData(decodedMsg, acc)

	default:
		return nil, 0, fmt.Errorf("unhandled operation type: %d", opType)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode message of operation type: %d: %w", opType, err)
	}

	addressBytes, err := acc.Address.MarshalText()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal address: %w", err)
	}

	_, minimalFees, err := network.GetNodeInfo()
	if err != nil {
		minimalFees = "0"
	}

	data.Description = description
	data.Fees = strconv.FormatUint(fees, 10)
	data.MinFees = minimalFees
	data.WalletAddress = string(addressBytes)
	data.Nickname = acc.Nickname
	data.OperationType = int(opType)
	data.AllowFeeEdition = *params.AllowFeeEdition
	data.ChainID = chainID

	promptRequest := prompt.PromptRequest{
		Action: walletapp.Sign,
		Data:   data,
	}

	return &promptRequest, fees, nil
}

func getCallSCPromptData(
	decodedMsg []byte,
	acc *account.Account,
) (PromptRequestSignData, error) {
	msg, err := callsc.DecodeMessage(decodedMsg)
	if err != nil {
		return PromptRequestSignData{}, err
	}

	return PromptRequestSignData{
		Coins:    strconv.FormatUint(msg.Coins, 10),
		Address:  msg.Address,
		Function: msg.Function,
	}, nil
}

func getExecuteSCPromptData(
	decodedMsg []byte,
	acc *account.Account,
) (PromptRequestSignData, error) {
	msg, err := executesc.DecodeMessage(decodedMsg)
	if err != nil {
		return PromptRequestSignData{}, err
	}

	return PromptRequestSignData{
		MaxCoins: strconv.FormatUint(msg.MaxCoins, 10),
	}, nil
}

func getRollPromptData(
	decodedMsg []byte,
	acc *account.Account,
) (PromptRequestSignData, error) {
	msg, err := sendoperation.RollDecodeMessage(decodedMsg)
	if err != nil {
		return PromptRequestSignData{}, err
	}

	return PromptRequestSignData{
		RollCount: msg.RollCount,
		Coins:     strconv.FormatUint(msg.RollCount*RollPrice, 10),
	}, nil
}

func (w *walletSign) getTransactionPromptData(
	decodedMsg []byte,
	acc *account.Account,
) (PromptRequestSignData, error) {
	msg, err := transaction.DecodeMessage(decodedMsg)
	if err != nil {
		return PromptRequestSignData{}, err
	}

	var recipientNickname string

	recipientAcc, err := w.prompterApp.App().Wallet.GetAccountFromAddress(msg.RecipientAddress)
	if err != nil {
		recipientNickname = ""
	} else {
		recipientNickname = recipientAcc.Nickname
	}

	return PromptRequestSignData{
		RecipientAddress:  msg.RecipientAddress,
		RecipientNickname: recipientNickname,
		Amount:            strconv.FormatUint(msg.Amount, 10),
	}, nil
}

// privateKeyFromCache return the private key from the cache or an error.
func (w *walletSign) privateKeyFromCache(
	acc *account.Account,
	correlationID *memguard.LockedBuffer,
) (*memguard.LockedBuffer, error) {
	key := CacheKey(correlationID.Bytes())

	value, err := w.gc.Get(key)
	if err != nil {
		if err.Error() == gcache.KeyNotFoundError.Error() {
			return nil, fmt.Errorf("%w: %w", utils.ErrCorrelationIDNotFound, err)
		}

		return nil, fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	if value == nil {
		return nil, fmt.Errorf("%w: %s", utils.ErrCache, "value is nil")
	}

	byteValue, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("%w: %s", utils.ErrCache, "value is not a byte array")
	}

	cacheValue := make([]byte, ed25519.PrivateKeySize)
	copy(cacheValue, byteValue)

	cipheredPrivateKey := memguard.NewBufferFromBytes(cacheValue)

	privateKey, err := Xor(cipheredPrivateKey, correlationID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", utils.ErrCache, err)
	}

	cipheredPrivateKey.Destroy()
	correlationID.Destroy()

	return privateKey, nil
}

func Xor(bufferA, bufferB *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	a := bufferA.Bytes()
	b := bufferB.Bytes()

	if len(a) != len(b) {
		return nil, fmt.Errorf("length of two arrays must be same, %d and %d", len(a), len(b))
	}
	result := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}

	return memguard.NewBufferFromBytes(result), nil
}

func CacheKey(correlationID models.CorrelationID) [32]byte {
	return blake3.Sum256(correlationID)
}

func expirationDuration() time.Duration {
	fromEnv := os.Getenv("BATCH_EXPIRATION_TIME")

	if fromEnv == "" {
		return defaultExpirationTime
	}

	duration, err := time.ParseDuration(fromEnv)
	if err != nil {
		return defaultExpirationTime
	}

	return duration
}

func generateCorrelationID() (*memguard.LockedBuffer, error) {
	rand := cryptorand.Reader

	correlationId := make([]byte, ed25519.PrivateKeySize)
	if _, err := io.ReadFull(rand, correlationId); err != nil {
		return nil, err
	}

	return memguard.NewBufferFromBytes(correlationId), nil
}
