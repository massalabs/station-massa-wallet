package wallet

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/massalabs/station-massa-wallet/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func verifyStatusCode(t *testing.T, resp *httptest.ResponseRecorder, statusCode int) {
	// Log body to simplify failure analysis.
	body := new(strings.Builder)
	_, _ = io.Copy(body, resp.Result().Body)

	assert.Equal(t, statusCode, resp.Result().StatusCode, fmt.Sprintf("the returned body is: %s", strings.TrimSpace(body.String())))
}

// cleanupTestData cleans up wallet created file.
func cleanupTestData(nicknames []string) error {
	for _, name := range nicknames {
		w, err := wallet.Load(name)
		if err != nil {
			return err
		}
		err = w.DeleteFile()
		if err != nil {
			return err
		}
	}

	return nil
}
