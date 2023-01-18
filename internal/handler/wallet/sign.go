package wallet

import (
	"crypto/ed25519"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"lukechampine.com/blake3"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"

	"github.com/btcsuite/btcutil/base58"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/password"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

// NewSign instantiates a sign Handler
// The "classical" way is not possible because we need to pass to the handler a password.PasswordAsker.
func NewSign(pwdPrompt password.Asker) operations.RestWalletSignOperationHandler {
	return &walletSign{pwdPrompt: pwdPrompt}
}

type walletSign struct {
	pwdPrompt password.Asker
}

// Handle handles a sign request.
func (s *walletSign) Handle(params operations.RestWalletSignOperationParams) middleware.Responder {

	wlt, resp := loadWallet(params.Nickname)
	if resp != nil {
		return resp
	}

	resp = unprotectWalletAskingPassword(wlt, s.pwdPrompt, params.Nickname)
	if resp != nil {
		return resp
	}

	pubKey := wlt.KeyPair.PublicKey
	privKey := wlt.KeyPair.PrivateKey

	digest, resp := digestOperationAndPubKey(params.Body.Operation, pubKey)
	if resp != nil {
		return resp
	}

	signature := ed25519.Sign(privKey, digest[:])

	return operations.NewRestWalletSignOperationOK().WithPayload(
		&models.Signature{
			PublicKey: "P" + base58.CheckEncode(pubKey, wallet.Base58Version),
			Signature: signature,
		})
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

// unprotectWalletAskingPassword asks for a password and unprotects the wallet.
func unprotectWalletAskingPassword(wallet *wallet.Wallet, prompter password.Asker, nickname string) middleware.Responder {
	clearPassword, err := prompter.Ask(nickname)
	if err != nil {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorCanceledAction,
				Message: errorCanceledAction,
			})
	}

	if len(clearPassword) == 0 {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorPasswordEmptyExecuteFct,
				Message: errorPasswordEmptyExecuteFct,
			})
	}

	err = wallet.Unprotect(clearPassword)
	if err != nil {
		return operations.NewRestWalletSignOperationInternalServerError().WithPayload(
			&models.Error{
				Code:    errorWrongPassword,
				Message: "Error : cannot uncipher the wallet: " + err.Error(),
			})
	}

	return nil
}

// digestOperationAndPubKey prepares the digest for signature.
func digestOperationAndPubKey(operation *strfmt.Base64, publicKey []byte) ([32]byte, middleware.Responder) {
	// reads operation to sign
	op, err := operation.MarshalText()
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
