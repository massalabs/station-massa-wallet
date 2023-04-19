package wallet

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func Test_exportFileWallet_handler(t *testing.T) {
	api, _, _, err := MockAPI()
	if err != nil {
		panic(err)
	}

	handler, exist := api.HandlerFor("get", "/rest/wallet/export/file/{nickname}")
	if !exist {
		panic("Endpoint doesn't exist")
	}

	// test empty configuration first.
	t.Run("Export file of unknown wallet", func(t *testing.T) {
		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/rest/wallet/export/file/%s", "nobody"), "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		verifyStatusCode(t, resp, 404)
	})

	// test with one wallet configuration.
	t.Run("Export file of created wallet", func(t *testing.T) {
		nickname := "trololol"
		password := "zePassword"
		_, err = wallet.Generate(nickname, password)
		if err != nil {
			t.Fatalf(err.Error())
		}

		resp, err := handleHTTPRequest(handler, "GET", fmt.Sprintf("/rest/wallet/export/file/%s", nickname), "")
		if err != nil {
			t.Fatalf("while serving HTTP request: %s", err)
		}

		verifyStatusCode(t, resp, 200)
		verifyHeader(t, resp, "Content-Type", "application/octet-stream")
		verifyHeader(t, resp, "Content-Disposition", fmt.Sprintf("attachment; filename=%q", "wallet_trololol.yml"))
		verifyBody(t, resp, nickname)

		err = cleanupTestData([]string{nickname})
		if err != nil {
			t.Fatalf("while cleaning up TestData: %s", err)
		}
	})
}

func verifyHeader(t *testing.T, resp *httptest.ResponseRecorder, headerName, headerValue string) {
	if resp.Header().Get(headerName) != headerValue {
		t.Fatalf("the header %s was: %s, want %s", headerName, resp.Header().Get(headerName), headerValue)
	}
}

func verifyBody(t *testing.T, resp *httptest.ResponseRecorder, nickname string) {
	body := resp.Body.String()
	if body == "" {
		t.Fatalf("the body was empty")
	}

	// check the first line
	if !strings.Contains(body, "Nickname: "+nickname) {
		t.Fatalf("the body doesn't contain the wallet nickname")
	}
}
