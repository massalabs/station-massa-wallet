package wallet

import (
	"fmt"
	"testing"
)

func Test_walletImport_Handle(t *testing.T) {
	t.Skip("Skipping Test_walletImport_Handle")

	api, _, _, err := MockAPI()
	if err != nil {
		panic(err)
	}
	createTestWallet(t, api, "precondition_wallet", `{"Nickname": "precondition_wallet", "Password": "1234"}`, 200)

	type want struct {
		statusCode int
	}

	privateKeyForTests := "S12XPyhXmGnx4hnx59mRUXPo6BDb18D6a7tA1xyAxAQPPFDUSNXA"

	privateKeyPromptKeyOK := PrivateKeyPrompt{PrivateKey: privateKeyForTests, Err: nil}
	privateKeyPromptKeyKO := PrivateKeyPrompt{PrivateKey: "S12ABCD", Err: nil}
	privateKeyPromptError := PrivateKeyPrompt{PrivateKey: "", Err: fmt.Errorf("Private key is required")}

	testsImportWallet := []struct {
		name         string
		nickname     string
		promptResult PrivateKeyPrompt
		want         want
	}{
		{"passing", "titi", privateKeyPromptKeyOK, want{statusCode: 204}},
		{"wrong privateKey format", "titi", privateKeyPromptKeyKO, want{statusCode: 500}},
		{"nickName Already taken", "precondition_wallet", privateKeyPromptKeyOK, want{statusCode: 500}},
		{"PrivateKey null", "titi", privateKeyPromptError, want{statusCode: 400}},
	}
	for _, tt := range testsImportWallet {
		t.Run(tt.name, func(t *testing.T) {
			// keyChan <- tt.promptResult // non blocking call as channel is buffered
			// pwdChan <- PasswordPrompt{Password: "1234", Err: nil}

			handler, exist := api.HandlerFor("post", "/rest/wallet/import/{nickname}")
			if !exist {
				panic("Endpoint doesn't exist")
			}

			resp, err := handleHTTPRequest(handler, "POST", fmt.Sprintf("/rest/wallet/import/%s", tt.nickname), "")
			if err != nil {
				t.Fatalf("while serving HTTP request: %s", err)
			}
			verifyStatusCode(t, resp, tt.want.statusCode)
		})
	}

	err = cleanupTestData([]string{"precondition_wallet", "titi"})
	if err != nil {
		t.Fatalf("while cleaning up TestData: %s", err)
	}
}
