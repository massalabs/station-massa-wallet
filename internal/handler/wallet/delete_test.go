package wallet

import (
	"fmt"
	"testing"
)

func Test_walletDelete_Handle(t *testing.T) {
	api, _, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	// To test wallet deletion, we need to create one first.
	createTestWallet(t, api, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

	// Only the status code is checked
	testsDelete := []struct {
		name           string
		walletNickname string
		statusCode     int
	}{
		{"passing", "precondition_wallet", 204},
		{"failing_wallet_does_not_exist", "wallet_does_not_exist", 500},
	}

	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			handler, exist := api.HandlerFor("DELETE", "/rest/wallet/{nickname}")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			resp, err := handleHTTPRequest(handler, "DELETE", fmt.Sprintf("/rest/wallet/%s", tt.walletNickname), "")
			if err != nil {
				t.Fatalf("while serving HTTP request: %s", err)
			}

			if resp.Result().StatusCode != tt.statusCode {
				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.statusCode)
			}
		})
	}
}
