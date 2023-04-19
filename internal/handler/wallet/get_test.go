package wallet

import (
	"encoding/json"
	"fmt"
	"testing"

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

		var wallets []models.Wallet
		err = json.Unmarshal(resp.Body.Bytes(), &wallets)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Wallet: %s", err)
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

		var wallet []models.Wallet
		err = json.Unmarshal(resp.Body.Bytes(), &wallet)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Wallet: %s", err)
		}

		for idx, nickname := range nicknames {
			if wallet[idx].Nickname != nickname {
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
	api, _, _, err := MockAPI()
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
	t.Run("Passed_list", func(t *testing.T) {
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

		err = cleanupTestData([]string{nickname})
		if err != nil {
			t.Fatalf("while cleaning up TestData: %s", err)
		}
	})
}
