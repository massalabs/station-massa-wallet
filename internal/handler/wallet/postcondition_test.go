package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/massalabs/thyra-plugin-massa-wallet/api/server/models"
)

// Clean up test data by listing all created wallets with tests and deleting them.
func cleanupTestData(t *testing.T) {
	t.Log("\n................... Start cleanupTestData ...................\n")
	// Configure the API server for the GET wallets operation
	apiGet, err := configureAPIServerGet()
	if err != nil {
		panic(err)
	}

	// Configure the API server for the DELETE wallet operation
	apiDelete, err := configureAPIServerDelete()
	if err != nil {
		panic(err)
	}

	// Get the handler for the GET wallets operation
	handlerGet, exist := apiGet.HandlerFor("get", "/rest/wallet")
	if !exist {
		t.Fatalf("Endpoint doesn't exist")
	}

	// Create a GET request for the wallets
	httpRequest, err := http.NewRequest("GET", "/rest/wallet", strings.NewReader(""))
	if err != nil {
		t.Fatalf(err.Error())
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	handlerGet.ServeHTTP(resp, httpRequest)
	if resp.Result().StatusCode != 200 {
		return
	}

	// Unmarshal the list of wallets from the response body
	var wallets []models.Wallet
	err = json.Unmarshal(resp.Body.Bytes(), &wallets)
	if err != nil {
		t.Fatalf("impossible to hydrate models.Wallet: %s", err)
	}

	// Iterate through the wallets and delete them
	t.Log("\n................... loop through all wallets ...................\n")
	for _, walletToDelete := range wallets {
		// Get the handler for the DELETE wallet operation
		handlerDelete, exist := apiDelete.HandlerFor("DELETE", "/rest/wallet/{nickname}")
		if !exist {
			t.Fatalf("Endpoint doesn't exist")
		}

		// Create a DELETE request for the wallet with the specified nickname
		t.Logf("\n................... Removing Wallet %v ...................\n", *walletToDelete.Nickname)
		httpRequest, err := http.NewRequest("DELETE", fmt.Sprintf("/rest/wallet/%s", *walletToDelete.Nickname), strings.NewReader(""))
		if err != nil {
			t.Fatalf(err.Error())
		}

		httpRequest.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handlerDelete.ServeHTTP(resp, httpRequest)
	}
	t.Log("\n................... cleanupTestData complete ...................\n")

}
