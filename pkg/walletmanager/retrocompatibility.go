package walletmanager

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/massalabs/station/pkg/logger"
)

// MigrateWallet moves the wallet from the old location (executable file) to the new one (user config).
func (w *Wallet) MigrateWallet() error {
	oldPath, err := GetWorkDir()
	if err != nil {
		return fmt.Errorf("reading config directory '%s': %w", oldPath, err)
	}

	files, err := os.ReadDir(oldPath)
	if err != nil {
		return fmt.Errorf("reading working directory '%s': %w", oldPath, err)
	}

	for _, f := range files {
		fileName := f.Name()
		oldFilePath := path.Join(oldPath, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			newFilePath := path.Join(w.WalletPath, fileName)

			// Skip if new file path exists
			if _, err := os.Stat(newFilePath); err == nil {
				continue
			}

			logger.Infof("Migrating wallet from", oldFilePath, "to", newFilePath)

			err = os.Rename(oldFilePath, newFilePath)
			if err != nil {
				logger.Errorf("moving account file from '%s' to '%s': %w", oldFilePath, newFilePath, err)
			}
		}
	}

	return nil
}

func GetWorkDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("getting executable path: %w", err)
	}

	if runtime.GOOS == "darwin" {
		// On macOS, the executable is in a subdirectory of the working directory.
		// We need to go up 4 levels to get the working directory.
		// wallet-plugin.app/Contents/MacOS/wallet-plugin
		return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(ex)))), nil
	}

	dir := filepath.Dir(ex)

	// Helpful when developing:
	// when running `go run`, the executable is in a temporary directory.
	if strings.Contains(dir, "go-build") {
		return ".", nil
	}

	return filepath.Dir(ex), nil
}
