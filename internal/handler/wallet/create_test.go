package wallet

import (
	"encoding/json"
	"testing"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/restapi/operations"
)

func Test_walletCreate_Handle(t *testing.T) {
	api, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{"passing", `{"Nickname": "toto", "Password": "1234"}`, 200},
		{"without Password", `{"Nickname": "toto"}`, 422},
		{"without Nickname", `{"Password": "1234"}`, 422},
		{"without Password and Nickname", `{}`, 422},
	}
	for _, tt := range tests {
		createTestWallet(t, api, tt.name, tt.body, tt.statusCode)
	}

	err = cleanupTestData([]string{"toto"})
	if err != nil {
		t.Fatalf("while cleaning up TestData: %s", err)
	}
}

// createTestWallet tests the creation of a wallet.
func createTestWallet(t *testing.T, api *operations.MassaWalletAPI, name string, inputBody string, statusCode int) {
	t.Run(name, func(t *testing.T) {
		resp, err := processHTTPRequest(api, "POST", "/rest/wallet", inputBody)
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		if resp.Result().StatusCode != statusCode {
			t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, statusCode)
		}

		if resp.Result().StatusCode != 200 {
			return
		}

		var wallet models.Wallet
		err = json.Unmarshal(resp.Body.Bytes(), &wallet)
		if err != nil {
			t.Fatalf("impossible to hydrate models.Wallet: %s", err)
		}

		var body operations.RestWalletCreateBody
		err = json.Unmarshal([]byte(inputBody), &body)
		if err != nil {
			t.Fatalf("impossible to hydrate operations.RestWalletCreateBody: %s", err)
		}

		if wallet.Nickname != body.Nickname {
			t.Fatalf("the wallet nickname was: %s, want %s", wallet.Nickname, `toto`)
		}
	})
}
