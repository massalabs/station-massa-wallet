package wallet

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Get the number of accounts in the wallet
func (w *Wallet) GetAccountCount() int {
	return len(w.AllAccounts())
}

func ClearAccounts(t *testing.T, walletPath string) {
	files, err := os.ReadDir(walletPath)
	assert.NoError(t, err)

	for _, f := range files {
		fileName := f.Name()
		filePath := path.Join(walletPath, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			os.Remove(filePath)
		}
	}

	assert.NoError(t, err)
}
