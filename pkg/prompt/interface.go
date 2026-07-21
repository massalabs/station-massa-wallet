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
	TryLock() bool
	Unlock()
	SelectBackupFilepath(nickname string) (string, error)
}

type PromptLocker struct {
	mutex     sync.Mutex
	PromptApp *walletapp.WalletApp
}

// TryLock atomically starts a prompt session. It returns true if the caller
// acquired the lock, or false if another prompt session is already active.
// The check-and-set is performed under a single mutex acquisition so that
// concurrent callers cannot both observe the prompter as free and proceed.
func (w *PromptLocker) TryLock() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.PromptApp.IsListening {
		return false
	}

	w.PromptApp.IsListening = true

	return true
}

func (w *PromptLocker) Unlock() {
	w.mutex.Lock()
	w.PromptApp.IsListening = false
	w.mutex.Unlock()
}

func (w *PromptLocker) App() *walletapp.WalletApp {
	return w.PromptApp
}
