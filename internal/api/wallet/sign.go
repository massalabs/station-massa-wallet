package wallet

import (
	"crypto/ed25519"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"lukechampine.com/blake3"

	"github.com/massalabs/thyra-plugin-massa-core/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-core/api/server/restapi/operations"

	"github.com/massalabs/thyra-plugin-massa-core/pkg/base58"
	"github.com/massalabs/thyra-plugin-massa-core/pkg/gui"
	"github.com/massalabs/thyra-plugin-massa-core/pkg/wallet"
)

//nolint:nolintlint,ireturn
func NewSign(walletStorage *sync.Map, app *fyne.App) operations.RestWalletSignOperationHandler {
	return &walletSign{walletStorage: walletStorage, app: app}
}

type walletSign struct {
	walletStorage *sync.Map
	app           *fyne.App
}

//nolint:nolintlint,ireturn,funlen
func (s *walletSign) Handle(params operations.RestWalletSignOperationParams) middleware.Responder {

	// retrieves key pair using wallet's nickname.
	if len(params.Nickname) == 0 {
		return operations.NewRestWalletSignOperationBadRequest().WithPayload(
			&models.Error{
				Code:    errorCodeWalletSignOperationEmptyNickname,
				Message: "Error: nickname field is mandatory.",
			})
	}

	wlt, err := wallet.Load(params.Nickname)
	if err != nil {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeGetWallet,
				Message: "Error cannot load wallet : " + err.Error(),
			})
	}

	clearPassword, err := gui.AskPassword(params.Nickname, s.app)
	if err != nil {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    ErrorCodeWalletCanceledAction,
				Message: ErrorCodeWalletCanceledAction,
			})
	}

	if len(clearPassword) == 0 {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    ErrorCodeWalletPasswordEmptyExecuteFct,
				Message: ErrorCodeWalletPasswordEmptyExecuteFct,
			})

	}

	err = wlt.Unprotect(clearPassword, 0)

	if err != nil {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletWrongPassword,
				Message: "Error : cannot uncipher the wallet : " + err.Error(),
			})
	}

	pubKey := wlt.KeyPairs[0].PublicKey
	privKey := wlt.KeyPairs[0].PrivateKey

	// reads operation to sign
	op, err := params.Body.Operation.MarshalText()
	if err != nil {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCodeWalletSignOperationRead,
				Message: "Error: while reading operation.",
			})
	}

	// signs operation

	digest := blake3.Sum256(append(pubKey, op...))

	signature := ed25519.Sign(privKey, digest[:])

	// format public key
	pubKeyB58VC := "P" + base58.VersionedCheckEncode(pubKey, 0)

	return operations.NewRestWalletSignOperationOK().WithPayload(
		&models.Signature{
			PublicKey: pubKeyB58VC,
			Signature: signature,
		})
}
