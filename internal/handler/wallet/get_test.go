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
		resp, err := processHTTPRequest(api, "GET", "/rest/wallet", "")
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
		resp, err := processHTTPRequest(api, "GET", "/rest/wallet", "")
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

	handler, exist := api.HandlerFor("get", "/rest/wallet/{nickname}")
	if !exist {
		panic("Endpoint doesn't exist")
	}

	// test empty configuration first.
	t.Run("Get unknown wallet", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/rest/wallet/%s", "nobody"), "")
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

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/rest/wallet/%s", nickname), "")
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

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/rest/wallet/%s?ciphered=false", nickname), "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		verifyStatusCode(t, resp, 200)
		verifyBodyWallet(t, resp, nickname)
		verifyPublicKeyIsPresent(t, resp, nickname)

		result := <-testResult

		if !result.Success {
			t.Fatalf("Expected success, got error")
		}

		msg := "Unprotect Success"
		if result.Data != msg {
			t.Fatalf(fmt.Sprintf("Expected error message to be %s, got %s", msg, result.Data))
		}

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

	fmt.Println("debug start")
	fmt.Println(body)
	fmt.Println("debug stop")

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
