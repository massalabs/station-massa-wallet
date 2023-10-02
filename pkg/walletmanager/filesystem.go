package walletmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

const (
	directoryName             = "massa-station-wallet"
	FileModeUserReadWriteOnly = 0o600
)

// Path returns the path where the account yaml file are stored.
func Path() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting user config directory: %w", err)
	}

	path := filepath.Join(configDir, directoryName)

	// create the directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("creating account directory '%s': %w", path, err)
		}
	}

	return path, nil
}

func (w *Wallet) accountPath(nickname string) (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, fmt.Sprintf("wallet_%s.yaml", nickname)), nil
}

func (w *Wallet) Persist(acc account.Account) error {
	filePath, err := w.accountPath(acc.Nickname)
	if err != nil {
		return err
	}

	data, err := acc.Marshal()
	if err != nil {
		return fmt.Errorf("marshaling account: %w", err)
	}

	err = os.WriteFile(filePath, data, FileModeUserReadWriteOnly)
	if err != nil {
		return fmt.Errorf("writing wallet to '%s: %w", filePath, err)
	}

	return nil
}

func (w *Wallet) Load(filePath string) (*account.Account, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading wallet from '%s': %w", filePath, err)
	}

	acc := account.NewEmpty()
	err = acc.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling account: %w", err)
	}

	return &acc, nil
}

func (w *Wallet) deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("deleting file '%s': %w", filePath, err)
	}

	return nil
}
