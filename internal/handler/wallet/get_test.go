package wallet

import (
	"encoding/json"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
)

func Test_walletGet_Handle(t *testing.T) {
	api, _, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	// test empty configuration first.
	t.Run("Passed_list_empty", func(t *testing.T) {
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

	// test with one wallet configuration.
	t.Run("Passed_list_with_wallets", func(t *testing.T) {
		// let's create the wallet first.
		createTestWallet(t, api, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

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

		if wallet[0].Nickname != "precondition_wallet" {
			t.Fatalf("the wallet nickname was: %s, want %s", wallet[0].Nickname, "precondition_wallet")
		}
	})

	err = cleanupTestData([]string{"precondition_wallet"})
	if err != nil {
		t.Fatalf("while cleaning up TestData: %s", err)
	}
}
