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
	accounts                *sync.Map // Mapping from nickname to account
	InvalidAccountNicknames []string  // List of invalid account nicknames
	WalletPath              string
	mutex                   sync.Mutex
}

// New creates a new wallet instance.
// If walletPath is empty, it will use the default wallet path.
func New(walletPath string) (*Wallet, error) {
	wallet := &Wallet{
		accounts: &sync.Map{},
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

	// Clear the accounts map
	w.accounts.Range(func(key interface{}, value interface{}) bool {
		w.accounts.Delete(key)
		return true
	})

	// Clear invalid accounts
	w.InvalidAccountNicknames = []string{}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, f := range files {
		fileName := f.Name()
		filePath := path.Join(w.WalletPath, fileName)

		if strings.HasPrefix(fileName, "wallet_") && strings.HasSuffix(fileName, ".yaml") {
			nickname := w.nicknameFromFilePath(filePath)

			_, found := w.accounts.Load(nickname)
			if found {
				continue
			}

			acc, err := w.Load(filePath)
			if err != nil {
				logger.Warnf("invalid account found: %s; %v", nickname, err)
				w.InvalidAccountNicknames = append(w.InvalidAccountNicknames, nickname)

				continue
			}

			err = w.AddAccount(acc, false, false)
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
func (w *Wallet) AddAccount(acc *account.Account, persist bool, force bool) error {
	if acc == nil {
		return fmt.Errorf("account is nil")
	}

	// Validate nickname uniqueness
	err := w.NicknameIsUnique(acc.Nickname)
	if err != nil && !force {
		return err
	}

	// Validate unique private key
	err = w.AddressIsUnique(acc.Address)
	if err != nil && !force {
		return err
	}

	if persist {
		err = w.Persist(*acc)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrPersistingAccount, err)
		}
	}

	w.accounts.Store(acc.Nickname, acc)

	return nil
}

// GenerateAccount generates a new account and adds it to the wallet.
// It returns the generated account.
// It destroys the password.
func (w *Wallet) GenerateAccount(password *memguard.LockedBuffer, nickname string) (*account.Account, error) {
	acc, err := account.Generate(password, nickname)
	if err != nil {
		return nil, fmt.Errorf("generating account: %w", err)
	}

	err = w.AddAccount(acc, true, false)
	if err != nil {
		return nil, fmt.Errorf("adding account in GenerateAccount: %w", err)
	}

	return acc, nil
}

// Get an account from the wallet by nickname
func (w *Wallet) GetAccount(nickname string) (*account.Account, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if acc, found := w.accounts.Load(nickname); found {
		return acc.(*account.Account), nil
	}

	accountPath, err := w.AccountPath(nickname)
	if err != nil {
		return nil, fmt.Errorf("getting account path: %w", err)
	}

	acc, err := w.Load(accountPath)
	if err != nil {
		logger.Errorf("loading account: %v", err)
		return nil, AccountNotFoundError
	}

	err = w.AddAccount(acc, false, false)
	if err != nil {
		return nil, fmt.Errorf("adding account in GetAccount: %w", err)
	}

	return acc, nil
}

func (w *Wallet) GetAccountFromAddress(needle string) (*account.Account, error) {
	var errInRange error
	var accFound *account.Account

	w.accounts.Range(func(_, value interface{}) bool {
		acc, ok := value.(*account.Account)
		if ok {
			address, err := acc.Address.MarshalText()
			if err != nil {
				errInRange = err
				return false
			}

			if needle == string(address) {
				accFound = acc
				return false
			}
		}

		return true
	})

	if errInRange != nil {
		return nil, errInRange
	}

	if accFound == nil {
		return nil, AccountNotFoundError
	}

	return accFound, nil
}

// Delete an account from the wallet
func (w *Wallet) DeleteAccount(nickname string) error {
	if _, found := w.accounts.Load(nickname); !found {
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

	w.accounts.Delete(nickname)

	return nil
}

func (w *Wallet) AllAccounts() []*account.Account {
	var accounts []*account.Account

	w.accounts.Range(func(_, value interface{}) bool {
		acc, ok := value.(*account.Account)
		if ok {
			accounts = append(accounts, acc)
		}

		return true
	})

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Nickname < accounts[j].Nickname
	})

	return accounts
}
