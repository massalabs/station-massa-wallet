package walletmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/massalabs/station-massa-wallet/pkg/types"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
)

const (
	directoryName             = "massa-station-wallet"
	FileModeUserReadWriteOnly = 0o600
)

var ErrUnmarshalAccount = errors.New("unmarshaling account")

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

func (w *Wallet) AccountPath(nickname string) (string, error) {
	return filepath.Join(w.WalletPath, Filename(nickname)), nil
}

// filename returns the wallet filename based on the given nickname.
func Filename(nickname string) string {
	return fmt.Sprintf("wallet_%s.yaml", nickname)
}

func (w *Wallet) nicknameFromFilePath(filePath string) string {
	_, nicknameFromFileName := filepath.Split(filePath)
	nicknameFromFileName = strings.TrimPrefix(nicknameFromFileName, "wallet_")

	return strings.TrimSuffix(nicknameFromFileName, ".yaml")
}

func (w *Wallet) Persist(acc account.Account) error {
	filePath, err := w.AccountPath(acc.Nickname)
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
		return nil, fmt.Errorf("%w: %w", ErrUnmarshalAccount, err)
	}

	// Nickname in the file is optional, if it's not set, we use the filename
	if acc.Nickname == "" {
		acc.Nickname = w.nicknameFromFilePath(filePath)
	}

	if !account.NicknameIsValid(acc.Nickname) {
		return nil, fmt.Errorf("%w: '%s'", account.ErrInvalidNickname, acc.Nickname) // TODO: add unit test
	}

	// Address in the file is optional, if it's not set, we use the public key
	if acc.Address.Object.Data == nil {
		acc.Address = *types.NewAddressFromPublicKey(&acc.PublicKey)
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
