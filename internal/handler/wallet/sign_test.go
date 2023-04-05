package wallet

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/api/server/models"
)

type want struct {
	statusCode int
}

var (
	PasswordPromptOK    PasswordPrompt = PasswordPrompt{Password: "1234", Err: nil}
	PasswordPromptKO    PasswordPrompt = PasswordPrompt{Password: "4321", Err: nil}
	PasswordPromptError PasswordPrompt = PasswordPrompt{Password: "1234", Err: errors.New("Error while getting password PasswordPrompt")}
)

type TestSign struct {
	name         string
	nickname     string
	body         string
	promptResult PasswordPrompt
	want         want
}

func Test_walletSign_Handle(t *testing.T) {
	api, channel, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	// let's create a new wallet.
	createTestWallet(t, api, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

	testsSign := []TestSign{
		{"passing", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptOK, want{statusCode: 200}},
		{"wrong password", "precondition_wallet", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptKO, want{statusCode: 500}},
		{"wrong nickname", "titi", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptOK, want{statusCode: 500}},
		{"password prompt error", "titi", `{"operation":"MjIzM3QyNHQ="}`, PasswordPromptError, want{statusCode: 500}},
	}
	for _, tt := range testsSign {
		t.Run(tt.name, func(t *testing.T) {
			if nil != &tt.promptResult {
				channel <- tt.promptResult // non blocking call as channel is buffered
			}

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

func Test_walletSign_Handle_Batch(t *testing.T) {
	api, channel, _, err := MockAPI()
	if err != nil {
		panic(err)
	}
	createTestWallet(t, api, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

	testSignNewBatch := TestSign{
		"passing", "precondition_wallet", `{"operation":"MjIzM3QyNHQ=","batch":true}`, PasswordPromptOK, want{statusCode: 200},
	}

	channel <- testSignNewBatch.promptResult

	handler, exist := api.HandlerFor("post", "/rest/wallet/{nickname}/signOperation")
	if !exist {
		panic("Endpoint doesn't exist")
	}

	resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/rest/wallet/%s/signOperation", testSignNewBatch.nickname), testSignNewBatch.body)
	if err != nil {
		t.Fatalf("while serving HTTP request: %s", err)
	}

	verifyStatusCode(t, resp, testSignNewBatch.want.statusCode)

	var body models.Signature
	err = json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("while unmarshalling: %s", err)
	}

	correlationId := base64.StdEncoding.EncodeToString(body.CorrelationID)

	testSignBatchItem := TestSign{
		"passing", "precondition_wallet", fmt.Sprintf(`{"operation":"MjIzM3QyNHQ=","correlationId":"%s"}`, correlationId), PasswordPromptOK, want{statusCode: 200},
	}

	resp, err = handleHTTPRequest(handler, "POST", fmt.Sprintf("/rest/wallet/%s/signOperation", testSignBatchItem.nickname), testSignBatchItem.body)
	if err != nil {
		t.Fatalf("while serving HTTP request (2nd): %s", err)
	}
	verifyStatusCode(t, resp, testSignBatchItem.want.statusCode)
}
