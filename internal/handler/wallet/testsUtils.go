package wallet

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/prompt"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/stretchr/testify/assert"
)

const (
	ChainIDUnitTests = 77658366 // Buildnet
)

func createAccount(password, nickname string, t *testing.T, prompterApp prompt.WalletPrompterInterface) *account.Account {
	acc, err := account.Generate(memguard.NewBufferFromBytes([]byte(password)), nickname)
	assert.NoError(t, err)
	err = prompterApp.App().Wallet.AddAccount(acc, true, false)
	assert.NoError(t, err)

	return acc
}

func verifyStatusCode(t *testing.T, resp *httptest.ResponseRecorder, statusCode int) {
	// Log body to simplify failure analysis.
	body := new(strings.Builder)
	_, _ = io.Copy(body, resp.Result().Body)

	assert.Equal(t, statusCode, resp.Result().StatusCode, "the returned body is: %s", strings.TrimSpace(body.String()))
}
