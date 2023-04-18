package wallet

import (
	"bytes"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"
	"time"

	"github.com/bluele/gcache"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"lukechampine.com/blake3"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

const passwordExpirationTime = time.Second * 60 * 30

// NewSign instantiates a sign Handler
// The "classical" way is not possible because we need to pass to the handler a password.PasswordAsker.
func NewSign(prompterApp wallet.WalletPrompterInterface, gc gcache.Cache) operations.RestWalletSignOperationHandler {
	return &walletSign{gc: gc, prompterApp: prompterApp}
}

type walletSign struct {
	prompterApp wallet.WalletPrompterInterface
	gc          gcache.Cache
}

// Handle handles a sign request.
func (s *walletSign) Handle(params operations.RestWalletSignOperationParams) middleware.Responder {
	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	var correlationId models.CorrelationID
	if params.Body.CorrelationID != nil {
		correlationId, resp = handleWithCorrelationId(wlt, params, s.gc)
	} else {
		if !wlt.UnprotectWalletAskingPassword(s.prompterApp) {
			return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
				&models.Error{
					Code:    errorCanceledAction,
					Message: "Unable to unprotect wallet",
				})
		}
		if params.Body.Batch {
			correlationId, resp = handleBatch(wlt, params, s, s.gc)
		}
	}
	if resp != nil {
		return resp
	}

	_, signature, resp := sign(wlt, params)
	if resp != nil {
		return resp
	}

	return operations.NewRestWalletSignOperationOK().WithPayload(
		&models.Signature{
			PublicKey:     wlt.GetPupKey(),
			Signature:     signature,
			CorrelationID: correlationId,
		})
}

func sign(wlt *wallet.Wallet, params operations.RestWalletSignOperationParams) ([]byte, []byte, middleware.Responder) {
	pubKey := wlt.KeyPair.PublicKey
	privKey := wlt.KeyPair.PrivateKey

	digest, resp := digestOperationAndPubKey(params.Body.Operation, pubKey)
	if resp != nil {
		return nil, nil, resp
	}

	signature := ed25519.Sign(privKey, digest[:])
	return pubKey, signature, nil
}

func handleWithCorrelationId(wlt *wallet.Wallet, params operations.RestWalletSignOperationParams, gc gcache.Cache) (models.CorrelationID, middleware.Responder) {
	cacheKey := getCacheKey(params.Body.CorrelationID)

	value, err := gc.Get(cacheKey)
	if err != nil {
		return nil, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignOperationLoadCache,
				Message: "Error cannot get data from cache: " + err.Error(),
			})
	}

	// convert interface{} into byte[]
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignOperationLoadCache,
				Message: "Error cannot convert cache value: " + err.Error(),
			})
	}
	bytes := buf.Bytes()

	err = wlt.UnprotectFromCorrelationId(bytes, params.Body.CorrelationID)

	if err != nil {
		return nil, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignOperationLoadCache,
				Message: "Error cannot unprotect from cache: " + err.Error(),
			})
	}

	return params.Body.CorrelationID, nil
}

func getCacheKey(correlationId models.CorrelationID) [32]byte {
	return blake3.Sum256(correlationId)
}

func handleBatch(wlt *wallet.Wallet, params operations.RestWalletSignOperationParams, s *walletSign, gc gcache.Cache) (models.CorrelationID, middleware.Responder) {
	correlationId, err := generateCorrelationId()
	if err != nil {
		return nil, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignOperationGenerateCorrelationId,
				Message: "Error cannot generate correlation id: " + err.Error(),
			})
	}

	cacheKey := getCacheKey(correlationId)
	cacheValue, err := wallet.Xor(wlt.KeyPair.PrivateKey, correlationId)
	if err != nil {
		return nil, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignOperationGenerateCorrelationId,
				Message: "Error cannot XOR correlation id: " + err.Error(),
			})
	}
	err = gc.SetWithExpire(cacheKey, cacheValue, passwordExpirationTime)
	if err != nil {
		return nil, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignOperationGenerateCorrelationId,
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

// loadWallet loads a wallet from the file system or returns an error.
func loadWallet(nickname string) (*wallet.Wallet, middleware.Responder) {
	if len(nickname) == 0 {
		return nil, operations.NewRestWalletSignOperationBadRequest().WithPayload(
			&models.Error{
				Code:    errorSignOperationEmptyNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	w, err := wallet.Load(nickname)
	if err != nil {
		return nil, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorGetWallet,
				Message: "Error cannot load wallet: " + err.Error(),
			})
	}

	return w, nil
}

// digestOperationAndPubKey prepares the digest for signature.
func digestOperationAndPubKey(operation *strfmt.Base64, publicKey []byte) ([32]byte, middleware.Responder) {
	// reads operation to sign

	op, err := base64.StdEncoding.DecodeString(operation.String())
	if err != nil {
		return [32]byte{}, operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorSignOperationRead,
				Message: "Error: while reading operation.",
			})
	}

	// signs operation
	digest := blake3.Sum256(append(publicKey, op...))

	return digest, nil
}
