package wallet

import (
	"os"
	"path/filepath"

	"github.com/massalabs/thyra-plugin-wallet/pkg/wallet"
)

// cleanupTestData cleans up wallet created file.
func cleanupTestData(nicknames []string) error {

	// get the current working directory
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, name := range nicknames {
		fullPath := filepath.Join(path, wallet.Filename(name))

		err := os.Remove(fullPath)
		if err != nil {
			return err
		}
	}

	return nil
}
