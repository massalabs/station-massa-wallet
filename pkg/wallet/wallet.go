package wallet

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/awnumar/memguard"
	"github.com/massalabs/station-massa-wallet/pkg/wallet/account"
	"github.com/massalabs/station/pkg/logger"
)

var (
	ErrPersistingAccount = errors.New("error persisting account")
	AccountNotFoundError = errors.New("account not found")
)

type Wallet struct {
	accounts                map[string]*account.Account // Mapping from nickname to account
	InvalidAccountNicknames []string                    // List of invalid account nicknames
	WalletPath              string
	mutex                   sync.Mutex
}

func New(walletPath string) (*Wallet, error) {
	wallet := &Wallet{
		accounts: make(map[string]*account.Account),
	}

	if walletPath == "" {
		walletPath, err := Path()
		if err != nil {
			return nil, fmt.Errorf("getting account directory: %w", err)
		}
		wallet.WalletPath = walletPath
	} else {
		wallet.WalletPath = walletPath
	}

	err := wallet.MigrateWallet()
	if err != nil {
		logger.Errorf("migrating wallet: %s", err)
	}

	err = wallet.Discover()
	if err != nil {
		return nil, fmt.Errorf("discovering accounts: %s\n", err)
	}

	return wallet, nil
}

func (w *Wallet) Discover() error {
	files, err := os.ReadDir(w.WalletPath)
	if err != nil {
		return fmt.Errorf("reading accounts path: %w", err)
	}

	for _, f := range files {
		fileName := f.Name()
		filePath := path.Join(w.WalletPath, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			nickname := w.nicknameFromFilePath(filePath)
			if w.accounts[nickname] != nil {
				continue
			}

			acc, err := w.Load(filePath)
			if err != nil {
				logger.Warnf("invalid account found: %s", nickname)
				w.InvalidAccountNicknames = append(w.InvalidAccountNicknames, nickname)

				continue
			}

			err = w.AddAccount(acc, false)
			if err != nil {
				logger.Warnf("failed to add account: %s, %v", nickname, err)
				w.InvalidAccountNicknames = append(w.InvalidAccountNicknames, nickname)

				continue
			}
		}
	}

	return nil
}

// Add an account into the wallet
func (w *Wallet) AddAccount(acc *account.Account, persist bool) error {
	if acc == nil {
		return fmt.Errorf("account is nil")
	}

	// Validate nickname uniqueness
	err := w.NicknameIsUnique(acc.Nickname)
	if err != nil {
		return err
	}

	// Validate unique private key
	err = w.AddressIsUnique(acc.Address)
	if err != nil {
		return err
	}

	if persist {
		err = w.Persist(*acc)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrPersistingAccount, err)
		}
	}

	w.addAccount(acc)

	return nil
}

func (w *Wallet) addAccount(acc *account.Account) {
	w.mutex.Lock()
	w.accounts[acc.Nickname] = acc
	w.mutex.Unlock()
}

// GenerateAccount generates a new account and adds it to the wallet.
// It returns the generated account.
// It destroys the password.
func (w *Wallet) GenerateAccount(password *memguard.LockedBuffer, nickname string) (*account.Account, error) {
	acc, err := account.Generate(password, nickname)
	if err != nil {
		return nil, fmt.Errorf("generating account: %w", err)
	}

	err = w.AddAccount(acc, true)
	if err != nil {
		return nil, fmt.Errorf("adding account: %w", err)
	}

	return acc, nil
}

// Get an account from the wallet by nickname
func (w *Wallet) GetAccount(nickname string) (*account.Account, error) {
	if w.accounts[nickname] != nil {
		return w.accounts[nickname], nil
	}

	accountPath, err := w.AccountPath(nickname)
	if err != nil {
		return nil, fmt.Errorf("getting account path: %w", err)
	}

	acc, err := w.Load(accountPath)
	if err != nil {
		logger.Errorf("loading account: %s", err)
		return nil, AccountNotFoundError
	}

	err = w.AddAccount(acc, false)
	if err != nil {
		return nil, fmt.Errorf("adding account: %w", err)
	}

	return acc, nil
}

// Delete an account from the wallet
func (w *Wallet) DeleteAccount(nickname string) error {
	if w.accounts[nickname] == nil {
		return AccountNotFoundError
	}

	accountPath, err := w.AccountPath(nickname)
	if err != nil {
		return fmt.Errorf("getting account path: %w", err)
	}

	err = w.deleteFile(accountPath)
	if err != nil {
		return fmt.Errorf("deleting account file: %w", err)
	}

	delete(w.accounts, nickname)

	return nil
}

// Get the number of accounts in the wallet
func (w *Wallet) GetAccountCount() int {
	return len(w.accounts)
}

func (w *Wallet) AllAccounts() []account.Account {
	accounts := make([]account.Account, 0, len(w.accounts))

	for _, acc := range w.accounts {
		accounts = append(accounts, *acc)
	}

	sort.SliceStable(accounts, func(i, j int) bool {
		return accounts[i].Nickname < accounts[j].Nickname
	})

	return accounts
}
