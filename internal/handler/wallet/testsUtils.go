package wallet

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

func verifyStatusCode(t *testing.T, resp *httptest.ResponseRecorder, statusCode int) {
	if resp.Result().StatusCode != statusCode {
		// Log body to simplify failure analysis.
		body := new(strings.Builder)
		_, _ = io.Copy(body, resp.Result().Body)

		t.Logf("the returned body is: %s", strings.TrimSpace(body.String()))

		t.Fatalf("the status code was: %d, want %d", resp.Result().StatusCode, statusCode)
	}
}

// cleanupTestData cleans up wallet created file.
func cleanupTestData(nicknames []string) error {
	for _, name := range nicknames {
		err := wallet.Delete(name)
		if err != nil {
			return err
		}
	}

	return nil
}
