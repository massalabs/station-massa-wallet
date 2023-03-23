package wallet

import (
	"errors"
	"fmt"
	"testing"
)

func Test_walletDelete_Handle(t *testing.T) {
	api, channel, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	// To test wallet deletion, we need to create one first.
	createTestWallet(t, api, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

	ConfirmPromptOK := &PasswordPrompt{Password: "1234", Err: nil}
	ConfirmPromptKO := &PasswordPrompt{Password: "4321", Err: nil}
	ConfirmPromptError := &PasswordPrompt{Password: "1234", Err: errors.New("Canceled by user")}

	testsDelete := []struct {
		name           string
		walletNickname string
		statusCode     int
		promptResult   *PasswordPrompt
	}{
		{"passing", "precondition_wallet", 204, ConfirmPromptOK},
		{"wrong password", "precondition_wallet", 500, ConfirmPromptKO},
		{"canceled by user", "precondition_wallet", 500, ConfirmPromptError},
		{"wrong nickname", "wallet_does_not_exist", 500, nil},
	}

	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			if nil != tt.promptResult {
				channel <- *tt.promptResult // non blocking call as channel is buffered
			}

			handler, exist := api.HandlerFor("DELETE", "/rest/wallet/{nickname}")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			resp, err := handleHTTPRequest(handler, "DELETE", fmt.Sprintf("/rest/wallet/%s", tt.walletNickname), "")
			if err != nil {
				t.Fatalf("while serving HTTP request: %s", err)
			}

			verifyStatusCode(t, resp, tt.statusCode)
		})
	}
}
