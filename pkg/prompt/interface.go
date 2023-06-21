package prompt

import (
	"sync"

	walletapp "github.com/massalabs/station-massa-wallet/pkg/app"
)

// WalletPrompterInterface is used to mock the WalletPrompter struct in tests.
type WalletPrompterInterface interface {
	PromptRequest(req PromptRequest)
	EmitEvent(eventId string, data walletapp.EventData)
	App() *walletapp.WalletApp
	IsListening() bool
	Unlock()
	Lock()
	SelectBackupFilepath(nickname string) (string, error)
}

type PromptLocker struct {
	mutex     sync.Mutex
	PromptApp *walletapp.WalletApp
}

func (w *PromptLocker) Lock() {
	w.mutex.Lock()
	w.PromptApp.IsListening = true
	w.mutex.Unlock()
}

func (w *PromptLocker) Unlock() {
	w.mutex.Lock()
	w.PromptApp.IsListening = false
	w.mutex.Unlock()
}

func (w *PromptLocker) IsListening() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.PromptApp.IsListening
}

func (w *PromptLocker) App() *walletapp.WalletApp {
	return w.PromptApp
}
