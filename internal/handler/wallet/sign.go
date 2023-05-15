package wallet

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
	"lukechampine.com/blake3"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

const passwordExpirationTime = time.Second * 60 * 30

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
	wlt, resp := loadWalletToSign(params.Nickname)
	if resp != nil {
		return resp
	}

	var correlationId models.CorrelationID
	if params.Body.CorrelationID != nil {
		correlationId, resp = handleWithCorrelationId(wlt, params, s.gc)
	} else {
		promptData := &prompt.PromptRequestData{
			Msg:  fmt.Sprintf("Unprotect wallet %s", wlt.Nickname),
			Data: nil,
		}

		_, err := prompt.PromptPassword(s.prompterApp, wlt, walletapp.Sign, promptData)
		if err != nil {
			return operations.NewSignUnauthorized().WithPayload(
				&models.Error{
					Code:    errorCanceledAction,
					Message: "Unable to unprotect wallet",
				})
		}

		s.prompterApp.EmitEvent(walletapp.PasswordResultEvent,
			walletapp.EventData{Success: true, Data: "Unprotect Success"})

		if params.Body.Batch {
			correlationId, resp = handleBatch(wlt, params, s, s.gc)
		}
	}
	if resp != nil {
		return resp
	}

	signature, err := wlt.Sign(params.Body.Operation)
	if err != nil {
		return operations.NewSignInternalServerError().WithPayload(
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
	cacheKey := getCacheKey(params.Body.CorrelationID)

	value, err := gc.Get(cacheKey)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: "Error cannot get data from cache: " + err.Error(),
			})
	}

	// convert interface{} into byte[]
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: "Error cannot convert cache value: " + err.Error(),
			})
	}
	bytes := buf.Bytes()

	err = wlt.UnprotectFromCorrelationId(bytes, params.Body.CorrelationID)

	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignLoadCache,
				Message: "Error cannot unprotect from cache: " + err.Error(),
			})
	}

	return params.Body.CorrelationID, nil
}

func getCacheKey(correlationId models.CorrelationID) [32]byte {
	return blake3.Sum256(correlationId)
}

func handleBatch(wlt *wallet.Wallet, params operations.SignParams, s *walletSign, gc gcache.Cache) (models.CorrelationID, middleware.Responder) {
	correlationId, err := generateCorrelationId()
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignGenerateCorrelationId,
				Message: "Error cannot generate correlation id: " + err.Error(),
			})
	}

	cacheKey := getCacheKey(correlationId)
	cacheValue, err := wallet.Xor(wlt.KeyPair.PrivateKey, correlationId)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignGenerateCorrelationId,
				Message: "Error cannot XOR correlation id: " + err.Error(),
			})
	}
	err = gc.SetWithExpire(cacheKey, cacheValue, passwordExpirationTime)
	if err != nil {
		return nil, operations.NewSignInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignGenerateCorrelationId,
				Message: "Error set correlation id in cache: " + err.Error(),
			})
	}

	return correlationId, nil
}

func generateCorrelationId() (models.CorrelationID, error) {
	rand := cryptorand.Reader

	correlationId := make([]byte, 64) // 64 is the private key size, correlation id must have the same size
	if _, err := io.ReadFull(rand, correlationId); err != nil {
		return nil, err
	}

	return correlationId, nil
}

// loadWalletToSign loads a wallet from the file system or returns an error.
func loadWalletToSign(nickname string) (*wallet.Wallet, middleware.Responder) {
	w, err := wallet.Load(nickname)

	if err == nil {
		return w, nil
	}

	errorObj := &models.Error{
		Code:    errorGetWallets,
		Message: err.Error(),
	}

	if err.Error() == wallet.ErrorAccountNotFound(nickname).Error() {
		return nil, operations.NewSignNotFound().WithPayload(errorObj)
	} else {
		return nil, operations.NewSignBadRequest().WithPayload(errorObj)
	}
}
