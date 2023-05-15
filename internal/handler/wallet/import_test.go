package wallet

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/api/server/restapi/operations"
	walletapp "github.com/massalabs/thyra-plugin-wallet/pkg/app"
	"github.com/massalabs/thyra-plugin-wallet/pkg/prompt"
)

func importWallet(t *testing.T, api *operations.MassaWalletAPI) *httptest.ResponseRecorder {
	handler, exist := api.HandlerFor("PUT", "/api/accounts")
	if !exist {
		panic("Endpoint doesn't exist")
	}

	resp, err := handleHTTPRequest(handler, "PUT", "/api/accounts", "")
	if err != nil {
		t.Fatalf("while serving HTTP request: %s", err)
	}
	return resp
}

func Test_walletImport_Handle(t *testing.T) {
	api, prompterApp, resChan, err := MockAPI()
	if err != nil {
		panic(err)
	}

	t.Run("import wallet file", func(t *testing.T) {
		nickname := "walletToBeImported"

		walletFile := fmt.Sprintf(
			`Version: 0
Nickname: %s
Address: AU12rjXkQ1hzo5hQ9Vvd4wdckdevZWAmT458WkbthGCgLoQ1C8BkB
Salt: [137, 167, 117, 16, 181, 186, 226, 139, 151, 228, 119, 194, 80, 40, 37, 138]
Nonce: [141, 224, 29, 20, 143, 96, 92, 165, 6, 227, 180, 8]
CipheredData: [9, 32, 11, 154, 176, 82, 30, 97, 4, 142, 233, 194, 12, 192, 138, 250,
  143, 241, 64, 45, 112, 130, 104, 211, 132, 86, 153, 68, 34, 243, 232, 24, 230, 136,
  54, 140, 186, 101, 64, 0, 194, 238, 105, 240, 46, 162, 168, 168, 238, 177, 175,
  202, 9, 17, 5, 153, 159, 34, 189, 231, 34, 116, 129, 72, 222, 174, 149, 22, 7, 210,
  167, 30, 252, 241, 63, 60, 190, 199, 8, 180]
PublicKey: [164, 243, 44, 155, 204, 6, 20, 131, 218, 97, 32, 58, 224, 189, 41, 113,
  4, 133, 59, 3, 213, 78, 33, 118, 49, 207, 180, 77, 78, 128, 37, 131]
`, nickname)

		filePath := "importMe.yml"
		// Write wallet file
		data := []byte(walletFile)
		err = os.WriteFile(filePath, data, 0o644)
		if err != nil {
			t.Fatalf(err.Error())
		}

		testResult := make(chan walletapp.EventData)

		// Send filepath to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().WalletFileChan <- filePath
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := importWallet(t, api)

		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "Import Success")

		assertWallet(t, nickname)

		err = cleanupTestData([]string{nickname})
		if err != nil {
			t.Fatalf("while cleaning up TestData: %s", err)
		}

		os.Remove(filePath)
	})

	t.Run("import invalid wallet file, then cancel", func(t *testing.T) {
		walletFile := "InvalidWalet"

		filePath := "importMe.yml"
		// Write wallet file
		data := []byte(walletFile)
		err = os.WriteFile(filePath, data, 0o644)
		if err != nil {
			t.Fatalf(err.Error())
		}

		testResult := make(chan walletapp.EventData)

		// Send filepath to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			// Send invalid filename to prompter app and wait for result
			prompterApp.App().WalletFileChan <- "invalidFilename"
			failRes := <-resChan

			checkResultChannel(t, failRes, false, prompt.InvalidAccountFileErr)

			// Send invalid filename to prompter app and wait for result
			prompterApp.App().WalletFileChan <- filePath
			failRes = <-resChan

			checkResultChannel(t, failRes, false, prompt.AccountLoadErr+": unmarshalling file")

			// Send cancel to prompter app to unlock the handler
			prompterApp.App().CtrlChan <- walletapp.Cancel
		}(testResult)

		resp := importWallet(t, api)

		verifyStatusCode(t, resp, http.StatusUnauthorized)
	})

	t.Run("import private key", func(t *testing.T) {
		nickname := "walletToBeImported"
		privateKey := "S12XPyhXmGnx4hnx59mRUXPo6BDb18D6a7tA1xyAxAQPPFDUSNXA"
		password := "aGoodPassword"
		testResult := make(chan walletapp.EventData)

		// Send account credentials to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PrivateKeyChan <- walletapp.ImportFromPKey{PrivateKey: privateKey, Nickname: nickname, Password: password}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := importWallet(t, api)

		verifyStatusCode(t, resp, http.StatusOK)

		result := <-testResult

		checkResultChannel(t, result, true, "Import Success")

		assertWallet(t, nickname)

		err = cleanupTestData([]string{nickname})
		if err != nil {
			t.Fatalf("while cleaning up TestData: %s", err)
		}
	})

	t.Run("import invalid private key", func(t *testing.T) {
		nickname := "walletToBeImported"
		privateKey := "invalidPrivateKey"
		password := "aWrongPassword"
		testResult := make(chan walletapp.EventData)

		// Send filepath to prompter app and wait for result
		go func(res chan walletapp.EventData) {
			prompterApp.App().PrivateKeyChan <- walletapp.ImportFromPKey{PrivateKey: privateKey, Nickname: nickname, Password: password}
			// forward test result to test goroutine
			res <- (<-resChan)
		}(testResult)

		resp := importWallet(t, api)

		verifyStatusCode(t, resp, http.StatusUnauthorized)

		result := <-testResult

		checkResultChannel(t, result, false, prompt.ImportPrivateKeyErr+": decoding private key: invalid format: version and/or checksum bytes missing")
	})
}
