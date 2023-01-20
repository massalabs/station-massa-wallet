package wallet

import (
	"errors"
	"fmt"
	"testing"
)

func Test_walletSign_Handle(t *testing.T) {
	api, channel, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	// let's create a new wallet.
	createTestWallet(t, api, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

	type want struct {
		statusCode int
	}

	PasswordPromptOK := PasswordPrompt{Password: "1234", Err: nil}
	PasswordPromptKO := PasswordPrompt{Password: "4321", Err: nil}
	PasswordPromptError := PasswordPrompt{Password: "1234", Err: errors.New("Error while getting password PasswordPrompt")}

	testsSign := []struct {
		name         string
		nickname     string
		body         string
		promptResult PasswordPrompt
		want         want
	}{
		{"passing", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptOK, want{statusCode: 200}},
		{"wrong password", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptKO, want{statusCode: 500}},
		{"wrong nickname", "titi", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptOK, want{statusCode: 500}},
		{"password prompt error", "titi", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptError, want{statusCode: 500}},
	}
	for _, tt := range testsSign {
		t.Run(tt.name, func(t *testing.T) {
			channel <- tt.promptResult // non blocking call as channel is buffered

			handler, exist := api.HandlerFor("post", "/rest/wallet/{nickname}/signOperation")
			if !exist {
				panic("Endpoint doesn't exist")
			}

			resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/rest/wallet/%s/signOperation", tt.nickname), tt.body)
			if err != nil {
				t.Fatalf("while serving HTTP request: %s", err)
			}

			verifyStatusCode(t, resp, tt.want.statusCode)
		})
	}

	err = cleanupTestData([]string{"precondition_wallet"})
	if err != nil {
		t.Fatalf("while cleaning up TestData: %s", err)
	}
}
