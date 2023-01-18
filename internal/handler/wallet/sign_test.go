package wallet

import (
	"errors"
	"fmt"
	"io"
	"strings"
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

	testsSign := []struct {
		name         string
		nickname     string
		body         string
		promptResult PasswordPrompt
		want         want
	}{
		{"passing", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, PasswordPrompt{Password: "1234", Err: nil}, want{statusCode: 200}},
		{"wrong password", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, PasswordPrompt{Password: "4321", Err: nil}, want{statusCode: 500}},
		{"wrong nickname", "titi", `{"operation":"MjIzM3QyNHQ="}`, PasswordPrompt{Password: "1234", Err: nil}, want{statusCode: 500}},
		{"PasswordPrompt error", "titi", `{"operation":"MjIzM3QyNHQ="}`, PasswordPrompt{Password: "1234", Err: errors.New("Error while getting password PasswordPrompt")}, want{statusCode: 500}},
	}
	for _, tt := range testsSign {
		t.Run(tt.name, func(t *testing.T) {
			channel <- tt.promptResult // non blocking call as channel is buffered

			handler, exist := api.HandlerFor("post", "/rest/wallet/{nickname}/signOperation")
			if !exist {
				t.Fatalf("Endpoint doesn't exist")
			}

			resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/rest/wallet/%s/signOperation", tt.nickname), tt.body)
			if err != nil {
				t.Fatalf("while serving HTTP request: %s", err)
			}

			if resp.Result().StatusCode != tt.want.statusCode {
				// Log body to simplify failure analysis.
				body := new(strings.Builder)
				_, _ = io.Copy(body, resp.Result().Body)

				t.Logf("the returned body is: %s", strings.TrimSpace(body.String()))

				t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, tt.want.statusCode)
			}
		})
	}

	err = cleanupTestData([]string{"precondition_wallet"})
	if err != nil {
		t.Fatalf("while cleaning up TestData: %s", err)
	}
}
