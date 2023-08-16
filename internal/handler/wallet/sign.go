package wallet

import (
	"bytes"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
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
	"lukechampine.com/blake3"
)

const passwordExpirationTime = time.Second * 60 * 30

type PromptRequestData struct {
	OperationType string
	OperationID   uint64
	GasLimit      uint64
	Coins         uint64
	Address       string
	Function      string
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

// Handle handles a sign request.
func (s *walletSign) Handle(params operations.SignParams) middleware.Responder {
	// params.Nickname length is already checked by go swagger
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	op, err := base64.StdEncoding.DecodeString(params.Body.Operation.String())
	if err != nil {
		return operations.NewSignBadRequest().WithPayload(
			&models.Error{
				Code:    errorSignDecodeOperation,
				Message: "Error: while decoding operation.",
			})
	}

	decodedMsg, err := sendoperation.DecodeMessage64(params.Body.Operation.String())
	if err != nil {
		return operations.NewSignBadRequest().WithPayload(
			&models.Error{
				Code:    errorSignDecodeMessage,
				Message: "Error: while decoding message.",
			})
	}

	callSC, err := callsc.DecodeMessage(decodedMsg)
	if err != nil {
		fmt.Println("fail decoding message, for now we decode only CallSC, ", err)
	}
	var correlationId models.CorrelationID
	if params.Body.CorrelationID != nil {
		correlationId, resp = handleWithCorrelationId(wlt, params, s.gc)
	} else {

		var promptRequest prompt.PromptRequest
		if callSC != nil {
			promptRequest = prompt.PromptRequest{
				Action: walletapp.Sign,
				Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
				Data: PromptRequestData{
					OperationType: "Call SC",
					OperationID:   callSC.OperationID,
					GasLimit:      callSC.GasLimit,
					Coins:         callSC.Coins,
					Address:       callSC.Address,
					Function:      callSC.Function,
				},
			}
		} else {
			promptRequest = prompt.PromptRequest{
				Action: walletapp.Sign,
				Msg:    fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
				Data: PromptRequestData{
					OperationType: "Unknown",
				},
			}
		}

		_, err := prompt.WakeUpPrompt(s.prompterApp, promptRequest, wlt)
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

	if err != nil {
		return operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignDecodeOperation,
				Message: "Error: while decoding operation.",
			})
	}

	signature, err := wlt.Sign(op)
	if err != nil {
		return operations.NewSignBadRequest().WithPayload(
			&models.Error{
				Code:    errorSignRead,
				Message: "Error: while reading operation.",
			})
	}

	return operations.NewSignOK().WithPayload(
		&models.SignResponse{
			PublicKey:     wlt.GetPupKey(),
			Signature:     signature,
			CorrelationID: correlationId,
		})
}

func handleWithCorrelationId(wlt *wallet.Wallet, params operations.SignParams, gc gcache.Cache) (models.CorrelationID, middleware.Responder) {
	key := CacheKey(params.Body.CorrelationID)

	value, err := gc.Get(key)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: fmt.Sprintf("Error cannot get data from cache: %v", err.Error()),
			})
	}

	// convert interface{} into byte[]
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: fmt.Sprintf("Error cannot convert cache value: %v", err.Error()),
			})
	}
	bytes := buf.Bytes()

	err = wlt.UnprotectFromCorrelationId(bytes, params.Body.CorrelationID)

	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: fmt.Sprintf("Error cannot unprotect from cache: %v", err.Error()),
			})
	}

	return params.Body.CorrelationID, nil
}

func CacheKey(correlationId models.CorrelationID) [32]byte {
	return blake3.Sum256(correlationId)
}

func handleBatch(wlt *wallet.Wallet, params operations.SignParams, s *walletSign, gc gcache.Cache) (models.CorrelationID, middleware.Responder) {
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
