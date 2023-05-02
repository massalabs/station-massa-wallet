package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func deleteWallet(t *testing.T, api *operations.MassaWalletAPI, nickname string) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("DELETE", "/api/accounts/{nickname}")
	if !exist {
		t.Fatalf("Endpoint doesn't exist")
	}

	resp, err := handleHTTPRequest(handler, "DELETE", fmt.Sprintf("/api/accounts/%s", nickname), "")
	if err != nil {
		t.Fatalf("while serving HTTP request: %s", err)
	}
	return resp
}

func Test_walletDelete_Handle(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	if err != nil {
		t.Fatalf(err.Error())
	}

	nickname := "walletToDelete"
	password := "zePassword"
	_, err = wallet.Generate(nickname, password)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("invalid nickname", func(t *testing.T) {
		resp := deleteWallet(t, api, "toto")
		verifyStatusCode(t, resp, http.StatusNotFound)
	})

	t.Run("invalid password", func(t *testing.T) {
		resp := deleteWallet(t, api, nickname)

		verifyStatusCode(t, resp, http.StatusNoContent)

		prompterApp.App().PasswordChan <- "invalid password"

		result := <-resChan

		checkResultChannel(t, result, false, "error unprotecting wallet:opening the private key seal: cipher: message authentication failed")
	})

	t.Run("canceled by user", func(t *testing.T) {
		resp := deleteWallet(t, api, nickname)

		verifyStatusCode(t, resp, http.StatusNoContent)

		prompterApp.App().CtrlChan <- walletapp.Cancel

		_, err = wallet.Load(nickname)
		if err != nil {
			t.Fatalf("Wallet should not have been deleted: " + err.Error())
		}
	})

	t.Run("delete success", func(t *testing.T) {
		resp := deleteWallet(t, api, nickname)

		verifyStatusCode(t, resp, http.StatusNoContent)

		prompterApp.App().PasswordChan <- password

		result := <-resChan

		checkResultChannel(t, result, true, "Delete Success")

		_, err = wallet.Load(nickname)
		if err == nil {
			t.Fatalf("Wallet should have been deleted")
		}
	})
}
