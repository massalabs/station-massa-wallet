package wallet

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
)

func Test_getWallets_handler(t *testing.T) {
	api, _, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	// test empty configuration first.
	t.Run("Get empty list", func(t *testing.T) {
		resp, err := processHTTPRequest(api, "GET", "/api/accounts", "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		if resp.Result().StatusCode != 200 {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, 200)
		}

		var wallets []models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Account: %s", err)
		}

		if len(wallets) != 0 {
			t.Fatalf("the wallets list should be empty")
		}
	})

	// Create wallets
	nicknames := []string{"wallet1", "wallet2", "wallet3"}
	for _, nickname := range nicknames {
		password := "zePassword"
		_, err = wallet.Generate(nickname, password)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

	t.Run("Get multiple wallets", func(t *testing.T) {
		resp, err := processHTTPRequest(api, "GET", "/api/accounts", "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		if resp.Result().StatusCode != 200 {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, 200)
		}

		var wallet []models.Account
		err = json.Unmarshal(resp.Body.Bytes(), &wallet)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Account: %s", err)
		}

		for idx, nickname := range nicknames {
			if wallet[idx].Nickname != models.Nickname(nickname) {
				t.Fatalf("the wallet nickname was: %s, want %s", wallet[idx].Nickname, nickname)
			}
		}
	})

	err = cleanupTestData(nicknames)
	if err != nil {
		t.Fatalf("while cleaning up TestData: %s", err)
	}
}

func Test_getWallet_handler(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	if err != nil {
		panic(err)
	}

	handler, exist := api.HandlerFor("get", "/api/accounts/{nickname}")
	if !exist {
		panic("Endpoint doesn't exist")
	}

	// test empty configuration first.
	t.Run("Get unknown wallet", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s", "nobody"), "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		verifyStatusCode(t, resp, 404)
	})

	// test with one wallet configuration.
	t.Run("Passed_get_ciphered", func(t *testing.T) {
		nickname := "trololol"
		password := "zePassword"
		_, err = wallet.Generate(nickname, password)
		if err != nil {
			t.Fatalf(err.Error())
		}

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s", nickname), "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		verifyStatusCode(t, resp, 200)
		verifyBodyWallet(t, resp, nickname)
		verifyPublicKeyIsNotPresent(t, resp, nickname)

		err = cleanupTestData([]string{nickname})
		if err != nil {
			t.Fatalf("while cleaning up TestData: %s", err)
		}
	})

	// test with un-ciphered data.
	t.Run("Passed_get_un-ciphered", func(t *testing.T) {
		testResult := make(chan walletapp.EventData)
		nickname := "trololol"
		password := "zePassword"
		_, err = wallet.Generate(nickname, password)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Send password to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PasswordChan <- password
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/api/accounts/%s?ciphered=false", nickname), "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		verifyStatusCode(t, resp, 200)
		verifyBodyWallet(t, resp, nickname)
		verifyPublicKeyIsPresent(t, resp, nickname)

		result := <-testResult

		checkResultChannel(t, result, true, "Unprotect Success")

		err = cleanupTestData([]string{nickname})
		if err != nil {
			t.Fatalf("while cleaning up TestData: %s", err)
		}
	})
}

func verifyBodyWallet(t *testing.T, resp *httptest.ResponseRecorder, nickname string) {
	body := resp.Body.String()
	if body == "" {
		t.Fatalf("the body was empty")
	}

	// check the first line
	if !strings.Contains(body, "\"nickname\":\""+nickname+"\"") {
		t.Fatalf("the body doesn't contain the wallet nickname")
	}
}

func verifyPublicKeyIsPresent(t *testing.T, resp *httptest.ResponseRecorder, nickname string) {
	body := resp.Body.String()
	if body == "" {
		t.Fatalf("the body was empty")
	}

	// check the first line
	if !strings.Contains(body, "publicKey\":\"P") {
		t.Fatalf("the body doesn't contain the wallet public key")
	}
}

func verifyPublicKeyIsNotPresent(t *testing.T, resp *httptest.ResponseRecorder, nickname string) {
	body := resp.Body.String()
	if body == "" {
		t.Fatalf("the body was empty")
	}

	// check the first line
	if !strings.Contains(body, "publicKey\":\"\"") {
		t.Fatalf("the body contains the wallet public key")
	}
}
