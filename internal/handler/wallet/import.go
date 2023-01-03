package wallet

import (
	"github.com/go-openapi/runtime/middleware"
	"lukechampine.com/blake3"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"

	"github.com/btcsuite/btcutil/base58"
	"github.com/massalabs/thyra-plugin-massa-wallet/pkg/wallet"
)

// checkDuplicate checks that the same wallet doesn't already exist.
func checkDuplicate(nickname string) middleware.Responder {
	_, err := wallet.Load(nickname)
	if err == nil {
		return operations.NewRestWalletImportInternalServerError().WithPayload(
			&models.Error{
				Code:    errorAlreadyExists,
				Message: "Error: a wallet with the same nickname already exists.",
			})
	}

	return nil
}

// decodeWalletAttributes decodes all wallet attributes.
func decodeWalletAttributes(rawPrivateKey string, rawPublicKey string, rawSalt string, rawNonce string) (privKey []byte, pubKey []byte, salt [16]byte, nonce [12]byte, resp middleware.Responder) {
	resp = operations.NewRestWalletCreateUnprocessableEntity()

	privKey, _, err := base58.CheckDecode(rawPrivateKey)
	if err != nil {
		return
	}

	pubKey, _, err = base58.CheckDecode(rawPublicKey)
	if err != nil {
		return
	}

	unsizedSalt, _, err := base58.CheckDecode(rawSalt)
	if err != nil {
		return
	}

	copy(salt[:], unsizedSalt)

	unsizedNonce, _, err := base58.CheckDecode(rawNonce)
	if err != nil {
		return
	}

	copy(nonce[:], unsizedNonce)

	resp = nil
	return
}

// failure generates an internal server error based on given error.
func failure(err error) middleware.Responder {
	return operations.NewRestWalletCreateInternalServerError().WithPayload(
		&models.Error{
			Code:    errorImportNew,
			Message: err.Error(),
		})
}

// HandleImport handles a import request.
func HandleImport(params operations.RestWalletImportParams) middleware.Responder {
	resp := checkDuplicate(params.Body.Nickname)
	if resp != nil {
		return resp
	}

	privateKey, publicKey, salt, nonce, resp := decodeWalletAttributes(
		params.Body.KeyPair.PrivateKey, params.Body.KeyPair.PublicKey,
		params.Body.KeyPair.Salt,
		params.Body.KeyPair.Nonce)
	if resp != nil {
		return resp
	}

	wlt, err := wallet.New(
		params.Body.Nickname,
		blake3.Sum256(publicKey),
		privateKey, publicKey,
		salt, nonce)
	if err != nil {
		return failure(err)
	}

	err = wlt.Persist()
	if err != nil {
		return failure(err)
	}

	return operations.NewRestWalletImportNoContent()
}
